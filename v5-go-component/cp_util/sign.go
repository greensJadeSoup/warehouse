package cp_util

import (
	"crypto/md5"
	"fmt"
)

func EncodeSign(body, nonce, timestamp, appSecret string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(body + nonce + timestamp + appSecret)))
}
