package service

import (
	"cetool/pkg/configs"
	"cetool/pkg/utils/ssh"
	"context"
	"os"
	"path/filepath"
	"strings"
)

const sshPath = "/etc/ssh/sshd_config"

type SSHService struct{}

type ISSHService interface {
	Download(node configs.Node, filePath string, remotePath string) error
	Upload(node configs.Node, localPath string, remotePath string, permission string) error
	Run(node configs.Node, shell string) (string, error)
	UploadTools(node configs.Node) (string, error)
}

func NewISSHService() ISSHService {
	return &SSHService{}
}

func (s *SSHService) Download(node configs.Node, filePath string, remotePath string) error {
	con := ssh.Build(node)

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	err = con.CopyFromRemote(context.Background(), f, remotePath)
	return err
}

func (s *SSHService) Upload(node configs.Node, localPath string, remotePath string, permission string) error {
	con := ssh.Build(node)

	_, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	file, err := os.Open(localPath)
	defer file.Close()
	if err != nil {
		return err
	}
	err = con.Upload(context.Background(), file, remotePath, permission)

	return err
}

func (s *SSHService) Run(node configs.Node, shell string) (string, error) {
	con := ssh.Build(node)

	result, err := con.Run(shell)
	return result, err
}

func (s *SSHService) UploadTools(node configs.Node) (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	//fmt.Println(executable)

	//con := ssh.Build(node)
	run, err := s.Run(node, "ls -l "+executable)
	if err != nil && run == "" {
		return "", err
	}

	if strings.Contains(run, executable) && err == nil {
		return executable, nil
	}

	processName := filepath.Base(executable)
	run, err = s.Run(node, "which "+processName)
	if err != nil && run == "" {
		return "", err
	}
	if !strings.Contains(run, "no "+processName) && err == nil {
		return processName, nil
	}
	err = s.Upload(node, executable, executable, "0777")
	if err != nil {
		return "", err
	}

	return executable, nil
}
