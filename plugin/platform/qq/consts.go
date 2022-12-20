package qq

const ModuleName = "mirai_qq"

const (
	DeviceInfoFilePath = "./device.json"
	SessionTokenPath   = "./session.token"
)

// Login Methods
const (
	LoginByToken   = LoginMethod("token")
	LoginByQrCode  = LoginMethod("qrcode")
	LoginByAccount = LoginMethod("account")
)

const (
	ConfigKeyAccount  = "qq.account"
	ConfigKeyPassword = "qq.password"
)
