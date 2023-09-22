package ssh

import (
	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

type WebSSHProvider struct {
}

func (s *WebSSHProvider) Register(container framework.Container) framework.NewInstance {
	return NewWebSSHService
}

func (s *WebSSHProvider) Boot(container framework.Container) error {
	return nil
}

func (s *WebSSHProvider) IsDefer() bool {
	return true
}

func (s *WebSSHProvider) Params(container framework.Container) []any {
	return []any{container}
}

func (s *WebSSHProvider) Name() string{
	return contract.SSHKey

}
