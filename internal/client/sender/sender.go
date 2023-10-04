package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AntonPashechko/gophkeeper/internal/client/config"
	"github.com/AntonPashechko/gophkeeper/internal/encrypt"
	"github.com/AntonPashechko/gophkeeper/internal/models"
	"github.com/go-resty/resty/v2"
)

const (
	register = "api/user/register"
	login    = "api/user/login"
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
		return fmt.Errorf("cannot send register request: %w", err)
	}

	//Нужно разобрать заголовки и забрать токен
	if code := resp.StatusCode(); code != http.StatusOK {
		return fmt.Errorf("request processing failed, code: %d", code)
	}

	resp.Header().Get("Authorization")

	//Получение header c токеном
	m.token = resp.Header().Get("Authorization")
	if m.token == `` {
		return fmt.Errorf("authorization header is missing")
	}

	return nil
}

func (m *sender) Login(login, password string) error {

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

	url := strings.Join([]string{m.cfg.ServerEndpoint, login}, "/")

	resp, err := req.Post(url)
	if err != nil {
		return fmt.Errorf("cannot send login request: %w", err)
	}

	//Нужно разобрать заголовки и забрать токен
	if code := resp.StatusCode(); code != http.StatusOK {
		return fmt.Errorf("request processing failed, code: %d", code)
	}

	resp.Header().Get("Authorization")

	//Получение header c токеном
	m.token = resp.Header().Get("Authorization")
	if m.token == `` {
		return fmt.Errorf("authorization header is missing")
	}

	return nil
}
