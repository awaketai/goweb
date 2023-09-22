package ssh

import (
	"context"
	"sync"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
	"golang.org/x/crypto/ssh"
)

type WebSSHService struct {
	container framework.Container
	clients   map[string]*ssh.Client
	lock      *sync.RWMutex
}

func NewWebSSHService(parasm ...any) (any, error) {
	container := parasm[0].(framework.Container)
	clients := make(map[string]*ssh.Client)
	lock := &sync.RWMutex{}
	return &WebSSHService{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

func (w *WebSSHService) GetClient(options ...contract.SSHOptions) (*ssh.Client, error) {
	logService := w.container.MustMake(contract.LogKey).(contract.Log)
	config, err := GetBaseConfig(w.container)
	if err != nil {
		return nil, err
	}

	for _, opt := range options {
		if err := opt(w.container, config); err != nil {
			return nil, err
		}
	}
	key := config.UniqueKey()
	w.lock.RLock()
	if db, ok := w.clients[key]; ok {
		w.lock.RUnlock()
		return db, nil
	}
	w.lock.RUnlock()
	// instance
	w.lock.Lock()
	defer w.lock.Unlock()
	addr := config.Host + ":" + config.Port
	client, err := ssh.Dial(config.Network, addr, config.ClientConfig)
	if err != nil {
		logService.Error(context.Background(), "ssh dial err:", map[string]any{
			"err:":  err,
			"addr:": addr,
		})
		return nil, err
	}
	w.clients[key] = client

	return client, nil
}
