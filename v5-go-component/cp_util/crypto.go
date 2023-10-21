package cp_util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
)

// 获得Md5加密值
func Md5Encrypt(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 获得Sha1加密
func Sha1Encrypt(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// 获得Hmac-Sha1加密
func HmacEncryptSha1(key string, data []byte) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write(data)
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func HmacEncryptSha256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Aes加密
func AesEncrypt(str, key string) (string, error) {
	keyData := []byte(key)
	if len(keyData) != 16 {
		return "", errors.New("aes key 长度必须等于16")
	}

	block, err := aes.NewCipher(keyData)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData := []byte(str)
	//origData = PKCS5Padding(origData, blockSize)
	origData = zeroPadding(origData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, keyData[:blockSize])
	aesData := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// aesData := origData
	blockMode.CryptBlocks(aesData, origData)

	return base64.URLEncoding.EncodeToString(aesData), nil
}

// Aes解密
func AesDecrypt(str, key string) (string, error) {
	keyData := []byte(key)
	if len(keyData) != 16 {
		return "", errors.New("aes key 长度必须等于16")
	}

	aesData, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		aesData, err = base64.URLEncoding.DecodeString(str)
		if err != nil {
			return "", err
		}
	}

	block, err2 := aes.NewCipher(keyData)
	if err2 != nil {
		return "", err2
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyData[:blockSize])
	origData := make([]byte, len(aesData))
	// origData := aesData
	blockMode.CryptBlocks(origData, aesData)
	//origData = PKCS5UnPadding(origData)
	origData = zeroUnPadding(origData)
	origStr := strings.Replace(string(origData), "\n", "", -1)

	return origStr, nil
}

func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func zeroUnPadding(origData []byte) []byte {
	//	length := len(origData)
	//	unpadding := int(origData[length-1])
	//	return origData[:(length - unpadding)]

	index := bytes.IndexByte(origData, 0)
	return origData[0:index]
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
