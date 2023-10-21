package cp_util

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	v "warehouse/v5-go-component/cp_util/validation"
)

func RegexpMatch(text, regular string) bool {
	result, err := regexp.Match(regular, []byte(text))

	if err != nil {
		return false
	}

	return result
}

func HasCH(text string) bool {
	regular := `^[\\u4e00-\\u9fa5]$`
	valid := v.Validation{}
	return valid.Match(text, regexp.MustCompile(regular), "match").Ok
}

func IsValidUserName(text string) bool {
	regular := `^[a-zA-Z]\w{5,15}$`
	valid := v.Validation{}
	return valid.Match(text, regexp.MustCompile(regular), "match").Ok
}

func IsValidPassword(text string) bool {
	regular := `^[\@A-Za-z0-9\!\#\$\%\^\&\*\.\~]{6,16}$`
	valid := v.Validation{}
	return valid.Match(text, regexp.MustCompile(regular), "match").Ok
}

func IsMobile(text string) bool {
	valid := v.Validation{}
	return valid.Mobile(text, "mobile").Ok
}

func IsEmail(text string) bool {
	valid := v.Validation{}
	return valid.Email(text, "email").Ok
}

// CheckFunc 校验函数
// minLen: 最小长度
// maxLen: 最大长度
// str: 需要校验的字符串
// regexpStr: 正则表达式
// return error: 错误
func CheckFunc(minLen, maxLen int, str, regexpStr string) error {
	strLen := len(str)
	if strLen < minLen || strLen > maxLen {
		return fmt.Errorf("the length is invalid, %v-%v", minLen, maxLen)
	}

	// 判断正则表达式是否有误
	regCom, err := regexp.Compile(regexpStr)
	if err != nil {
		tmpStr := fmt.Sprintf("expression of regexp=%v is err: %v", regexpStr, err)
		return errors.New(tmpStr)
	}

	// 对 string 进行校验
	matchFlag := regCom.MatchString(str)
	if !matchFlag {
		tmpStr := fmt.Sprintf("params not match, is invalid")
		return errors.New(tmpStr)
	}

	return nil
}

// NumLetter 数字和字母
func NumLetter(minLen, maxLen int, str string) error {
	regexpStr := "^[a-zA-Z0-9]*$"

	return CheckFunc(minLen, maxLen, str, regexpStr)
}

// NumLetter 数字和字母
func NumLetterSymbol(minLen, maxLen int, str string) error {
	regexpStr := "^([a-zA-Z0-9]|[~!@#$|%^&*'><?.,\\/\"()])*$"

	return CheckFunc(minLen, maxLen, str, regexpStr)
}

// NumCheck 数字
func NumCheck(minLen, maxLen int, str string) error {
	regexpStr := "^[0-9]*$"

	return CheckFunc(minLen, maxLen, str, regexpStr)
}

// LetterCheck 英文字母
func LetterCheck(minLen, maxLen int, str string) error {
	regexpStr := "^[a-zA-Z]*$"

	return CheckFunc(minLen, maxLen, str, regexpStr)
}

// 中国大陆手机号码验证
func ChinaPhoneCheck(minLen, maxLen int, phone string) error {
	regexpStr := "^1[0-9]*$"
	return CheckFunc(minLen, maxLen, phone, regexpStr)
}

// 验证qq
func CheckQQ(minLen, maxLen int, str string) error {
	regexpStr := "^[1-9][0-9]*$"

	return CheckFunc(minLen, maxLen, str, regexpStr)
}

// PwdCheck 密码校验, 以数字和字母开头,包含下划线和扛
func PwdCheck(minLen, maxLen int, str string) error {
	regexpStr := "^[a-zA-Z0-9][a-zA-Z0-9_-]*$"

	return CheckFunc(minLen, maxLen, str, regexpStr)
}

func CheckRealName(realName string, min int, max int) bool {
	realName = strings.Trim(realName, "")
	count := utf8.RuneCountInString(realName)
	if count < min || count > max {
		return false
	}
	match_cn, err_cn := regexp.MatchString("^[\u4e00-\u9fa5]+([·•][\u4e00-\u9fa5]+)*$", realName)
	match_en, err_en := regexp.MatchString("^[a-zA-Z]+([\\s·•]?[a-zA-Z]+)+$", realName)
	if (!match_cn || err_cn != nil) && (!match_en || err_en != nil) {
		return false
	}
	return true
}
