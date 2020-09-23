package crypto

import (
	"crypto/aes"
	"fmt"
	"testing"
)

func Test_AESEncryptCBC(t *testing.T) {
	secretKey := "347f36057c5373fab0d69158f345bf8d"
	iv := "7ca646d6dbf731aa0af9e77b"
	origin := "Hello aes"

	sign, err := AESEncryptCBC([]byte(secretKey), []byte(iv), []byte(origin))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(sign)
}

func Test_AESDecryptCBC(t *testing.T) {
	sign := "ClnJjHTpMasT3FnELakVQQ=="
	secretKey := "347f36057c5373fab0d69158f345bf8d"
	iv := "7ca646d6dbf731aa0af9e77b963f1d41"
	value, err := AESDecryptCBC([]byte(secretKey), []byte(iv)[:aes.BlockSize], sign)
	fmt.Println(value, err)

}
