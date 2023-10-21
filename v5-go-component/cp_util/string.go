package cp_util

import (
"crypto/rand"
	"fmt"
	r "math/rand"
	"strings"
	"time"
"unsafe"
)



// 生成随机数字的字符串
func NewRandomNum(n int) string {
	return string(NewRandomBytes(n, "0123456789"))
}

// 生成随机字符串
// keywords 为可选参数，只取[0]，指定随机数取值范围
func NewRandomBytes(n int, keywords ...string) []byte {
	ks := ""
	if len(keywords) > 0 {
		ks = keywords[0]
	}
	alphabets := []byte(ks)

	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return bytes
}

func SubStr(str string, begin, length int) string {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

//高性能字符串转字节
func StringToByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//高性能字节转字符串
func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//Int64数组去重
func UniqueArrInt64(arr []int64) []int64 {
	newArr := make([]int64, 0)
	tempArr := make(map[int64]bool, len(newArr))
	for _, v := range arr {
		if tempArr[v] == false {
			tempArr[v] = true
			newArr = append(newArr, v)
		}
	}
	return newArr
}

//Int64数组分割字符串
func ArrayInt64ToString(items []int64) string {
	str := strings.Replace(strings.Trim(fmt.Sprint(items), "[]"), " ", ",", -1)
	return str
}

//下划线转大写驼峰
func ToUpperCamelCase(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.Title(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ToUpper(s[:1]) + s[1:]
	s = strings.ReplaceAll(s, "Id", "ID")
	return s
}

//下划线转小写驼峰
func ToLowerCamelCase(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Split(s, " ")
	for i := range words {
		if i == 0 {
			words[i] = strings.ToLower(words[i])
		} else {
			words[i] = strings.Title(words[i])
		}
	}
	s = strings.Join(words, "")
	s = strings.ReplaceAll(s, "Id", "ID")
	return s
}