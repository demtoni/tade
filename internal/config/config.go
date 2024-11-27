package config

import "os"

type Config struct {
	ServerAddr    string
	PathToDB      string
	SessionSecret string
	Domain        string
	YooShopID     string
	YooApiKey     string
}

func New() (*Config, error) {
	// TODO: Error-check config parameters
	return &Config{
		ServerAddr:    os.Getenv("SERVER_ADDR"),
		PathToDB:      os.Getenv("DB_PATH"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		Domain:        os.Getenv("DOMAIN_NAME"),
		YooApiKey:     os.Getenv("YOO_API_KEY"),
		YooShopID:     os.Getenv("YOO_SHOP_ID"),
	}, nil
}
