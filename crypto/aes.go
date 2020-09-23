package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

func AESEncryptCBC(key, iv []byte, origData []byte) (string, error) {
	if len(iv) < 16 {
		return "", fmt.Errorf("iv length must lte 16")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	iv = iv[:blockSize]
	origData = PKCS7Padding(origData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)

	return base64.StdEncoding.EncodeToString(crypted), nil

}

func AESDecryptCBC(key, iv []byte, crypted string) (string, error) {
	if len(iv) < 16 {
		return "", fmt.Errorf("iv length must lte 16")
	}

	decodeData, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	iv = iv[:blockSize]

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(decodeData))

	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("CryptBlocks %v", e)
		}
	}()

	blockMode.CryptBlocks(origData, decodeData)

	origData, err = UnPadding(origData)

	return string(origData), err
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func UnPadding(origData []byte) ([]byte, error) {
	if len(origData) == 0 {
		return origData, nil
	}
	length := len(origData)
	// 去掉最后一个字节 unPadding 次
	unPadding := int(origData[length-1])
	if len(origData) >= (length-unPadding) && (length-unPadding) >= 0 {
		return origData[:(length - unPadding)], nil
	}
	return origData, errors.New("UnPadding error, please check key")
}
