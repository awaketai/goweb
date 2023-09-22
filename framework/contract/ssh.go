package contract

import (
	"fmt"

	"github.com/awaketai/goweb/framework"
	"golang.org/x/crypto/ssh"
)

const SSHKey = "web:ssh"

type SSH interface{
	GetClient(option ...SSHOptions) (*ssh.Client,error)
}

type SSHOptions func(container framework.Container,config *SSHConfig) error

type SSHConfig struct{
	Network string
	Host string
	Port string
	*ssh.ClientConfig
}


func (s *SSHConfig) UniqueKey() string{
	
	return fmt.Sprintf("%v_%v_%v", s.Host,s.Port,s.User)
}
