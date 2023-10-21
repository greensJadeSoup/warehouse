package cp_util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	math_rand "math/rand"
	"time"
	"unsafe"
)

/*==================================算法1=====================================*/
func NewGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}



/*================================算法2=======================================*/

const letters1 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letters2 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letters3 = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letters4 = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = math_rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandStrUpper(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters1) {
			b[i] = letters1[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func RandStrUpperLower(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters2) {
			b[i] = letters2[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func RandStrUpperNum(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters3) {
			b[i] = letters3[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func RandStrUpperLowerNum(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters4) {
			b[i] = letters4[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}



