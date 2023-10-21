package cp_util

import (
	"os/exec"
	"runtime"
	"strings"
	"warehouse/v5-go-component/cp_log"
)

func RunInLinuxWithErr(cmd string) (string, error) {
	if runtime.GOOS != "linux" {
		return "", nil
	}

	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		cp_log.Error(err.Error())
		return "", err
	}

	return strings.TrimSpace(string(result)), err
}
