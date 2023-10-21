package cp_util

import (
	"os"
	"runtime"
	"warehouse/v5-go-component/cp_log"
)

const CANGBOSS_UUID = "42005b8c-b736-48ec-a67a-581ffce4f81a"

func CheckMoralCharacter()  {
	if runtime.GOOS != "linux" {
		return
	}

	result, err := RunInLinuxWithErr("echo c3ad23705cd742b9191f5cf629b504bd | sudo -S dmidecode -t system | grep UUID | awk '{print $2}'")
	if err != nil {
		cp_log.Info("RunInLinuxWithErr:" + err.Error())
		os.Exit(-100)
	} else if result != CANGBOSS_UUID {
		os.Exit(-250)
	} else {
		cp_log.Info("check success")
	}
}

