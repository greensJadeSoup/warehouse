package cp_util

import (
	"errors"
	"fmt"
	"strconv"
)

func Bytes2str(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

//保留小数后xx位
func RemainBit(f float64, bit int) (float64, error) {
	if bit == 1 {
		return strconv.ParseFloat(fmt.Sprintf("%.1f", f), 64)
	} else if bit == 2 {
		return strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
	} else if bit == 3 {
		return strconv.ParseFloat(fmt.Sprintf("%.3f", f), 64)
	}

	return -1, errors.New("保留小数转换失败")
}