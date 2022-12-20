package qq

import (
	"bytes"
	"os"
	"strings"
	"time"

	"github.com/TeemoKill/BBot/log"
	"github.com/TeemoKill/BBot/utils"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	miraiBinary "github.com/Mrs4s/MiraiGo/binary"
	miraiClient "github.com/Mrs4s/MiraiGo/client"
	"github.com/pkg/errors"
	"github.com/tuotoo/qrcode"
)

type LoginMethod string

func (m *MiraiQQ) login() (err error) {
	logger := log.CurrentModuleLogger()
	logger.WithField("protocol", miraiClient.SystemDeviceInfo.Protocol).
		Info("开始尝试登录并同步消息...")

	fileExist, _ := utils.FileExist(SessionTokenPath)
	if fileExist {
		var token []byte
		token, err = os.ReadFile(SessionTokenPath)
		if err != nil {
			logger.WithError(err).
				WithField("session_token_path", SessionTokenPath).
				Errorf("读取会话令牌失败，可能是会话令牌已损坏，尝试删除令牌重新登录")
			return err
		}
		if m.qqClient.Uin != 0 {
			r := miraiBinary.NewReader(token)
			sessionUin := r.ReadInt64()
			if sessionUin != m.qqClient.Uin {
				logger.Warnf("QQ号(%d)与会话缓存内的QQ号(%d)不符，将清除会话缓存", m.qqClient.Uin, sessionUin)
				_ = m.clearToken()
				err = m.normalLogin()
				if err != nil {
					logger.WithError(err).
						Errorf("mirai normalLogin error")
					return err
				}
			}
		}
		err = m.qqClient.TokenLogin(token)
		if err != nil {
			_ = m.clearToken()
			logger.WithError(err).Warnf("恢复会话失败, 尝试使用正常流程登录。")
			time.Sleep(time.Second)
			m.resetMiraiClient(
				m.Cfg.GetInt64(ConfigKeyAccount),
				m.Cfg.GetString(ConfigKeyPassword),
			)
		} else {
			err = m.saveToken(SessionTokenPath)
			if err != nil {
				logger.WithError(err).
					Errorf("qq module saveToken error")
				return err
			}
			logger.Debug("恢复会话成功")
			return nil
		}
	} else {
		err = m.normalLogin()
		if err != nil {
			logger.WithError(err).
				Errorf("qq module normalLogin error")
			return err
		}
	}

	return nil
}

func (m *MiraiQQ) normalLogin() (err error) {
	logger := log.CurrentModuleLogger()

	if m.qqClient.Uin == 0 {
		logger.Info("未指定账号密码，请扫码登陆")
		err = m.qrcodeLogin()
		if err != nil {
			logger.WithError(err).
				Errorf("login failed")
			return err
		}
	} else {
		logger.Info("使用帐号密码登陆")
		err = m.commonLogin()
		if err != nil {
			logger.WithError(err).
				Errorf("login failed")
			return err
		}
	}
	err = m.saveToken(SessionTokenPath)
	if err != nil {
		logger.WithError(err).
			Errorf("save token error")
		return err
	}

	return nil
}

func (m *MiraiQQ) qrcodeLogin() (err error) {
	logger := log.CurrentModuleLogger()

	qrCodeResp, err := m.qqClient.FetchQRCode()
	if err != nil {
		logger.WithError(err).Errorf("mirai client fetch qr code error")
		return err
	}
	qrMatrix, err := qrcode.Decode(bytes.NewReader(qrCodeResp.ImageData))
	if err != nil {
		logger.WithError(err).Errorf("qrcode Decode error")
		return err
	}
	_ = os.WriteFile("qrcode.png", qrCodeResp.ImageData, 0o644)
	defer func() { _ = os.Remove("qrcode.png") }()
	if m.qqClient.Uin != 0 {
		logger.Infof("请使用账号 %d 登录手机QQ扫描二维码 (qrcode.png) : ", m.qqClient.Uin)
	} else {
		logger.Infof("请使用手机QQ扫描二维码 (qrcode.png) : ")
	}
	time.Sleep(time.Second)
	qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite, qrcodeTerminal.QRCodeRecoveryLevels.Low).Get(qrMatrix.Content).Print()
	qrCodeStatus, err := m.qqClient.QueryQRCodeStatus(qrCodeResp.Sig)
	if err != nil {
		logger.WithError(err).Errorf("mirai QueryQRCodeStatus error")
		return err
	}

	prevState := qrCodeStatus.State
	for {
		time.Sleep(time.Second)
		qrCodeStatus, _ = m.qqClient.QueryQRCodeStatus(qrCodeResp.Sig)
		if qrCodeStatus == nil {
			continue
		}
		if prevState == qrCodeStatus.State {
			continue
		}
		prevState = qrCodeStatus.State
		switch qrCodeStatus.State {
		case miraiClient.QRCodeCanceled:
			logger.Fatalf("扫码被用户取消.")
		case miraiClient.QRCodeTimeout:
			logger.Fatalf("二维码过期")
		case miraiClient.QRCodeWaitingForConfirm:
			logger.Infof("扫码成功, 请在手机端确认登录.")
		case miraiClient.QRCodeConfirmed:
			var resp *miraiClient.LoginResponse
			resp, err = m.qqClient.QRCodeLogin(qrCodeStatus.LoginInfo)
			if err != nil {
				logger.WithError(err).Errorf("mirai QRCodeLogin error")
				return err
			}
			return m.processLoginResp(resp)
		case miraiClient.QRCodeImageFetch, miraiClient.QRCodeWaitingForScan:
			// ignore
		}
	}
}

func (m *MiraiQQ) commonLogin() (err error) {
	logger := log.CurrentModuleLogger()

	loginResp, err := m.qqClient.Login()
	if err != nil {
		logger.WithError(err).Errorf("mirai client login error")
		return err
	}

	err = m.processLoginResp(loginResp)
	return err
}

func (m *MiraiQQ) processLoginResp(loginResp *miraiClient.LoginResponse) (err error) {
	logger := log.CurrentModuleLogger()

	// process login response
	for loginResp != nil {
		if loginResp.Success {
			return nil
		}

		// login not success yet
		var text string
		switch loginResp.Error {
		case miraiClient.SliderNeededError:
			logger.Warnf("登录需要滑条验证码, 请使用手机QQ扫描二维码以继续登录.")
			m.resetMiraiClient(0, "")
			return m.qrcodeLogin()
		case miraiClient.NeedCaptcha:
			logger.Warnf("登录需要验证码.")
			_ = os.WriteFile("captcha.jpg", loginResp.CaptchaImage, 0o644)
			logger.Warnf("请输入验证码 (captcha.jpg)： (Enter 提交)")
			text, err = utils.ReadLine()
			if err != nil {
				logger.WithError(err).Errorf("error during ReadLine")
				return err
			}
			_ = os.Remove("captcha.jpg")
			loginResp, err = m.qqClient.SubmitCaptcha(text, loginResp.CaptchaSign)
			if err != nil {
				logger.WithError(err).Errorf("mirai submit captcha error")
				return err
			}
			continue
		case miraiClient.SMSNeededError:
			logger.Warnf("账号已开启设备锁, 按 Enter 向手机 %s 发送短信验证码.", loginResp.SMSPhone)
			_, _ = utils.ReadLine()
			if !m.qqClient.RequestSMS() {
				logger.Warnf("发送验证码失败，可能是请求过于频繁.")
				return errors.WithStack(MiraiSMSRequestErr)
			}
			logger.Warn("请输入短信验证码： (Enter 提交)")
			text, err = utils.ReadLine()
			if err != nil {
				logger.WithError(err).Errorf("error during ReadLine")
				return err
			}
			loginResp, err = m.qqClient.SubmitSMS(text)
			if err != nil {
				logger.WithError(err).Errorf("mirai SubmitSMS error")
				return err
			}
			continue
		case miraiClient.SMSOrVerifyNeededError:
			logger.Warnf("账号已开启设备锁，请选择验证方式:")
			logger.Warnf("1. 向手机 %v 发送短信验证码", loginResp.SMSPhone)
			logger.Warnf("2. 使用手机QQ扫码验证.")
			logger.Warn("请输入(1 - 2) (将在10秒后自动选择2)：")
			text, err = utils.ReadLineTimeout(time.Second*10, "2")
			if err != nil {
				logger.WithError(err).Errorf("error during ReadLineTimeout")
				return err
			}
			if strings.Contains(text, "1") {
				if !m.qqClient.RequestSMS() {
					logger.Warnf("发送验证码失败，可能是请求过于频繁.")
					return errors.WithStack(MiraiSMSRequestErr)
				}
				logger.Warn("请输入短信验证码： (Enter 提交)")
				text, err = utils.ReadLine()
				if err != nil {
					logger.WithError(err).Errorf("error during ReadLine")
					return err
				}
				loginResp, err = m.qqClient.SubmitSMS(text)
				if err != nil {
					logger.WithError(err).Errorf("mirai SubmitSMS error")
					return err
				}
				continue
			}
			fallthrough
		case miraiClient.UnsafeDeviceError:
			logger.Warnf("账号已开启设备锁，请前往 -> %s <- 验证后重启Bot.", loginResp.VerifyUrl)
			logger.Infof("按 Enter 或等待 5s 后继续....")
			_, _ = utils.ReadLineTimeout(time.Second*5, "")
			return MiraiUnsafeDeviceErr
		case miraiClient.OtherLoginError, miraiClient.UnknownLoginError, miraiClient.TooManySMSRequestError:
			msg := loginResp.ErrorMessage
			if strings.Contains(msg, "版本") {
				msg = "密码错误或账号被冻结"
			}
			if strings.Contains(msg, "冻结") {
				msg = "账号被冻结"
			}
			logger.Warnf("登录失败: %s", msg)
			logger.Infof("按 Enter 或等待 5s 后继续....")
			_, _ = utils.ReadLineTimeout(time.Second*5, "")
			return MiraiOtherLoginErr
		}
	}

	return MiraiNilLoginRespErr
}

func (m *MiraiQQ) saveToken(tokenPath string) (err error) {
	err = os.WriteFile(tokenPath, m.qqClient.GenToken(), 0o677)
	return err
}

func (m *MiraiQQ) clearToken() (err error) {
	return os.Remove(SessionTokenPath)
}
