package util

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"regexp"
	"runtime"
)

func Urandom() string {
	f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

//logs
func AddLog(err interface{}, v ...interface{}) {
	if _, ok := err.(error); ok {
		pc, _, line, _ := runtime.Caller(1)
		logs.Error("[Server] ", runtime.FuncForPC(pc).Name(), line, v, err)
	} else {
		logs.Info("[Server] ", err)
	}
}

func JudgeIp(ip string) (err error) {
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		err = fmt.Errorf("not validate IP address")
		AddLog(err)
		return
	}
	return
}
