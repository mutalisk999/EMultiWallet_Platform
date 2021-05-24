package authcode

import (
	"github.com/mojocn/base64Captcha"
)

type AuthCode struct {
	AuthCodeId    string
	Base64PicData string
}

func CreateAuthCode(height int, width int, codeLen int) *AuthCode {
	//var digitConfig = base64Captcha.ConfigDigit{
	//	Height:     height,
	//	Width:      width,
	//	MaxSkew:    maxSkew,
	//	DotCount:   dotCount,
	//	CaptchaLen: codeLen,
	//}

	var configC = base64Captcha.ConfigCharacter{
		Height:             height,
		Width:              width,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         codeLen,
	}


	codeId, digitCode := base64Captcha.GenerateCaptcha("", configC)
	base64String := base64Captcha.CaptchaWriteToBase64Encoding(digitCode)

	authCode := new(AuthCode)
	authCode.AuthCodeId = codeId
	authCode.Base64PicData = base64String

	return authCode
}

func VerifyAuthCode(authCodeId string, verifyValue string) bool {
	return base64Captcha.VerifyCaptcha(authCodeId, verifyValue)
}
