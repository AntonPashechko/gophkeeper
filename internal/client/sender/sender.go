package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AntonPashechko/gophkeeper/internal/client/config"
	"github.com/AntonPashechko/gophkeeper/internal/encrypt"
	"github.com/AntonPashechko/gophkeeper/internal/models"
	"github.com/go-resty/resty/v2"
)

const (
	register = "api/user/register"
)

// sender для взаимодействия клиента с сервером
type sender struct {
	cfg       *config.Config
	client    *resty.Client // клиент http
	encryptor *encrypt.Encryptor
	token     string
}

func NewSender(cfg *config.Config) sender {

	if !strings.HasPrefix(cfg.ServerEndpoint, "http") && !strings.HasPrefix(cfg.ServerEndpoint, "https") {
		cfg.ServerEndpoint = "http://" + cfg.ServerEndpoint
	}

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

	return nil
}

func (m *sender) Register(login, password string) error {

	req := m.client.R()

	authdto := new(bytes.Buffer)
	if err := json.NewEncoder(authdto).Encode(&models.AuthDTO{
		Login:    login,
		Password: password,
	}); err != nil {
		return fmt.Errorf("error encoding auth dto %w", err)
	}

	//Шифруем аутентификационные данные
	encryptbuf, err := m.encryptor.Encrypt(authdto.Bytes())
	if err != nil {
		return fmt.Errorf("cannot encrypt metrics: %w", err)
	}

	req.SetHeader("Content-Type", "application/json").
		SetBody(encryptbuf)

	url := strings.Join([]string{m.cfg.ServerEndpoint, register}, "/")

	resp, err := req.Post(url)
	if err != nil {
		return err
	}

	resp.Header()

	return nil
}
