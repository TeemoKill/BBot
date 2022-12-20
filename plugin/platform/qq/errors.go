package qq

import "errors"

var (
	MiraiReadDeviceInfoJsonErr = errors.New("mirai failed to read device info json")

	MiraiNilClientErr    = errors.New("mirai client is nil, maybe not initialized")
	MiraiNilLoginRespErr = errors.New("mirai login response is nil")

	MiraiSMSRequestErr   = errors.New("mirai sms request error")
	MiraiUnsafeDeviceErr = errors.New("mirai login device unsafe")
	MiraiOtherLoginErr   = errors.New("mirai login failed")
)
