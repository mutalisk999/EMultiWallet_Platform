package coin

import (
	"testing"
	"encoding/hex"
	"github.com/mutalisk999/bitcoin-lib/src/utility"
	"fmt"
)

func TestCoinVerifyTrx2(t *testing.T) {
	pubkeyUncompress, _ := hex.DecodeString("f7bbbb0a687190933eeae1d819b92e6d5d3bf2911c2e39ccb4d3a7e21c46c7a498503e6f8052ad535c4c5d47ae3310696fc8245baf5ada54e47977aec245a73f")
	signatureScriptBytes, _ := hex.DecodeString("3044022024cd4a0abd4b80d23cb06d3d6cb4ae6cc7cccd520a70fd839be11ec0b4326e05022017391ff4e3e9bc6e6167bec902f4fa95dec04cb5e8d615af76da28f6598d2221")
	hashBytes := utility.Sha256([]byte("FFSHKYKXII"))

	pubkeyCompress := make([]byte, 33, 33)
	if pubkeyUncompress[63]%2 == 0 {
		pubkeyCompress[0] = 0x2
	} else {
		pubkeyCompress[0] = 0x3
	}
	copy(pubkeyCompress[1:], pubkeyUncompress[0:32])

	verifyOk, err := CoinVerifyTrx2(pubkeyCompress, hashBytes, signatureScriptBytes)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !verifyOk {
		fmt.Println("verify signature error")
	}
	fmt.Println("OK")
}