package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type SnifferConfig struct {
	Interface string
}

type DBConfig struct {
	Host     string
	Port     uint16
	Username string
	Password string
	Database string
}

func load() error {
	return godotenv.Load()
}

func ProvideSnifferConfig() (SnifferConfig, error) {
	if err := load(); err != nil {
		return SnifferConfig{}, err
	}

	interfaceName := os.Getenv("INTERFACE")
	if interfaceName == "" {
		return SnifferConfig{}, errors.New("INTERFACE is not set")
	}

	return SnifferConfig{
		Interface: interfaceName,
	}, nil
}

func ProvideDBConfig() (DBConfig, error) {
	if err := load(); err != nil {
		return DBConfig{}, err
	}

	host := os.Getenv("HOST")
	if host == "" {
		return DBConfig{}, errors.New("HOST is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return DBConfig{}, errors.New("PORT is not set")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return DBConfig{}, err
	}

	username := os.Getenv("USERNAME")
	if username == "" {
		return DBConfig{}, errors.New("USERNAME is not set")
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		return DBConfig{}, errors.New("PASSWORD is not set")
	}

	database := os.Getenv("DATABASE")
	if password == "" {
		return DBConfig{}, errors.New("DATABASE is not set")
	}

	return DBConfig{
		Host:     host,
		Port:     uint16(portInt),
		Username: username,
		Password: password,
		Database: database,
	}, nil
}
