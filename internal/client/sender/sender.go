package sender

import (
	"fmt"

	"github.com/AntonPashechko/gophkeeper/internal/client/config"
	"github.com/AntonPashechko/gophkeeper/internal/encrypt"
	"github.com/go-resty/resty/v2"
)

// sender для взаимодействия клиента с сервером
type sender struct {
	cfg       *config.Config
	client    *resty.Client // клиент http
	encryptor *encrypt.Encryptor
	token     string
}

func NewSender(cfg *config.Config) sender {
	return sender{
		cfg:    cfg,
		client: resty.New(),
	}
}

func (m *sender) Init() error {

	encryptor, err := encrypt.NewEncryptor(m.cfg.CryptoKey)
	if err != nil {
		return fmt.Errorf("cannot create credentions encryptor: %w", err)
	}

	m.encryptor = encryptor

	if m.cfg.Login == `` {
		return fmt.Errorf("login is undefined")
	}
	if m.cfg.Password == `` {
		return fmt.Errorf("password is password")
	}

	return nil
}
