package cp_util

import (
	"encoding/csv"
	"os"
	"warehouse/v5-go-component/cp_error"
)

//records是一个二维数组，[第n行][第M列]
func LoadCsvCfg(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if reader == nil {
		return nil, cp_error.NewNormalError("NewReader return nil, file:", file)
	}

	return reader.ReadAll()
}
