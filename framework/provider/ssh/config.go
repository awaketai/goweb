package ssh

import (
	"context"
	"os"
	"time"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func GetBaseConfig(c framework.Container) (*contract.SSHConfig, error) {
	logService := c.MustMake(contract.LogKey).(contract.Log)
	config := &contract.SSHConfig{
		ClientConfig: &ssh.ClientConfig{
			Auth:            []ssh.AuthMethod{},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
	opt := WithConfigPath("ssh")
	err := opt(c, config)
	if err != nil {
		logService.Error(context.Background(), "parse config err:", map[string]any{
			"err:": err,
		})
		return nil, err
	}

	return config, nil
}

func WithConfigPath(path string) contract.SSHOptions {
	return func(container framework.Container, config *contract.SSHConfig) error {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		logService := container.MustMake(contract.LogKey).(contract.Log)
		conf := configService.GetStringMapString(path)
		// load ssh config
		if network, ok := conf["network"]; ok {
			config.Network = network
		}

		if host, ok := conf["host"]; ok {
			config.Host = host
		}
		if port, ok := conf["port"]; ok {
			config.Port = port
		}
		if username, ok := conf["username"]; ok {
			config.User = username
		}
		if password, ok := conf["password"]; ok {
			authPwd := ssh.Password(password)
			config.Auth = append(config.Auth, authPwd)
		}
		if rsaKey, ok := conf["rsa_key"]; ok {
			key, err := os.ReadFile(rsaKey)
			if err != nil {
				logService.Error(context.Background(), "read rsa_key err:", map[string]any{
					"key":  rsaKey,
					"path": path,
					"err:": err,
				})
				return err
			}
			singer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				logService.Error(context.Background(), "create rsa_key singer err:", map[string]any{
					"key":  rsaKey,
					"path": path,
					"err:": err,
				})
				return err
			}
			rsaKeyAuth := ssh.PublicKeys(singer)
			config.Auth = append(config.Auth, rsaKeyAuth)
		}
		if knowHosts, ok := conf["know_hosts"]; ok {
			hostKeyCallback, err := knownhosts.New(knowHosts)
			if err != nil {
				logService.Error(context.Background(), "knowhosts err:", map[string]any{
					"key":  knowHosts,
					"path": path,
					"err:": err,
				})
				return err
			}
			config.HostKeyCallback = hostKeyCallback
		}
		if timeout, ok := conf["timeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.Timeout = t
		}

		return nil
	}
}

func WithSSHConfig(f func(options *contract.SSHConfig)) contract.SSHOptions {
	return func(container framework.Container, config *contract.SSHConfig) error {
		f(config)
		return nil
	}
}
