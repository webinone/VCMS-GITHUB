package libs

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Encoding (text string) string {

	bytes := []byte(text)
	h := sha256.New()                       // new sha256 object
	h.Write(bytes)                          // data is now converted to hex
	code 	:= h.Sum(nil)                      // code is now the hex sum
	enc_string := hex.EncodeToString(code)     // converts hex to string

	return enc_string
}


