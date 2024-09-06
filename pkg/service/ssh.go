package service

import (
	"context"
	"inspect/pkg/configs"
	"inspect/pkg/utils/ssh"
	"os"
)

const sshPath = "/etc/ssh/sshd_config"

type SSHService struct{}

type ISSHService interface {
	Download(node configs.Node, filePath string, remotePath string) error
	Upload(node configs.Node, localPath string, remotePath string, permission string) error
	Run(node configs.Node, shell string) (string, error)
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
