package utils

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"inspect/pkg/global"
	"net"
	"strings"
)

func CheckError(detail string, err error) {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		global.LOG.Warn("请检查网络连接：", err)
		return
	}
	var exitError *ssh.ExitError
	ok := errors.As(err, &exitError)
	if !ok {
		global.LOG.Warnf("执行失败： %s", err.Error())
		return
	}
	if strings.Contains(err.Error(), "unable to authenticate") {
		global.LOG.Warn("身份验证失败：", err)
		return
	}
	if strings.Contains(err.Error(), "no such host") {
		global.LOG.Warn("请检查节点地址：", err)
		return
	}
	global.LOG.Warn("执行失败：", detail)
}
