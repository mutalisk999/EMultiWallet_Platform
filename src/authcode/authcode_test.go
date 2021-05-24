package authcode

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestAuthCode(t *testing.T) {
	authCode := CreateAuthCode(80, 240, 6)
	fmt.Println("authcode", *authCode)

	file, _ := os.Create("authcode.png")
	bytesPic, _ := base64.StdEncoding.DecodeString(strings.Split(authCode.Base64PicData, ",")[1])
	file.Write(bytesPic)
	file.Close()
}

func TestVerifyAuthCode(t *testing.T) {
	fmt.Println(VerifyAuthCode("JEHV1d8GMA2CHyoNd01K", "143290"))
}
