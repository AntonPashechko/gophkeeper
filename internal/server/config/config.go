package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/AntonPashechko/gophkeeper/internal/logger"
)

type Config struct {
	Endpoint    string
	DataBaseDNS string
	CryptoKey   string        //Путь до файла с приватным ключом сервера для расшифровывания данных
	JWTKey      []byte        //Ключ для создания/проверки jwt для авторизации
	JWTDuration time.Duration //Время действия jwt для авторизации
}

func Create() (*Config, error) {
	cfg := &Config{}

	var JWTKey, JWTDuration string

	flag.StringVar(&cfg.Endpoint, "a", "localhost:8033", "address and port to run server")
	flag.StringVar(&cfg.DataBaseDNS, "d", "", "db dns")
	flag.StringVar(&cfg.CryptoKey, "k", "private.rsa", "Server private key path")
	flag.StringVar(&JWTKey, "jwtk", "aL6HmkWp7D", "JWT key")
	flag.StringVar(&JWTDuration, "jwtt", "60m", "JWT duration")

	flag.Parse()

	/*Но если заданы в окружении - берем оттуда*/
	if addr, exist := os.LookupEnv("RUN_ADDRESS"); exist {
		cfg.Endpoint = addr
	}

	if dns, exist := os.LookupEnv("DATABASE_URI"); exist {
		logger.Info("DATABASE_URI env: %s", dns)
		cfg.DataBaseDNS = dns
	}

	if cfg.DataBaseDNS == `` {
		return nil, fmt.Errorf("db dns is empty")
	}

	if key, exist := os.LookupEnv("JWT_KEY"); exist {
		JWTKey = key
	}

	if duration, exist := os.LookupEnv("JWT_DURATION"); exist {
		JWTDuration = duration
	}

	cfg.JWTKey = []byte(JWTKey)
	if duration, err := time.ParseDuration(JWTDuration); err != nil {
		return nil, fmt.Errorf("JWT DURATION: %w", err)
	} else {
		cfg.JWTDuration = duration
	}

	return cfg, nil
}
