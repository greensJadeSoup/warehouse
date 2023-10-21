package cp_util

import (
	"bufio"
	"fmt"
	"os"
	"warehouse/v5-go-component/cp_error"
)

func DirMidirWhenNotExist(dir string) error {
	_, e := os.Stat(dir)
	if e != nil {
		if os.IsNotExist(e) {
			if err := os.Mkdir(dir, 0666); err != nil {
				return cp_error.NewSysError(err)
			}
		}
	}

	return nil
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func CreateFile(filename string, data string) (int, error) {
	file, err := os.Create(filename)
	if err != nil {
		return 0, cp_error.NewSysError(err)
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	defer write.Flush()

	return fmt.Fprintln(write, string(data))
}

