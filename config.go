package main

import "strconv"

type Config struct {
	DownloadPath string
	ClipsAmount  int
}

func GetConfig() *Config {
	return nil
}

func SaveConfig() {

}

func SetConfig(path string, amount string) (*Config, error) {
	intAmount, err := strconv.Atoi(amount)
	if err != nil {
		return nil, err
	}
	cfg := Config{
		DownloadPath: path,
		ClipsAmount:  intAmount,
	}
	return &cfg, nil
}
