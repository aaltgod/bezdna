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

type ServerConfig struct {
	Host string
	Port uint16
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

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return DBConfig{}, errors.New("POSTGRES_HOST is not set")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		return DBConfig{}, errors.New("POSTGRES_PORT is not set")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return DBConfig{}, err
	}

	username := os.Getenv("POSTGRES_USERNAME")
	if username == "" {
		return DBConfig{}, errors.New("POSTGRES_USERNAME is not set")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return DBConfig{}, errors.New("POSTGRES_PASSWORD is not set")
	}

	database := os.Getenv("POSTGRES_DATABASE")
	if password == "" {
		return DBConfig{}, errors.New("POSTGRES_DATABASE is not set")
	}

	return DBConfig{
		Host:     host,
		Port:     uint16(portInt),
		Username: username,
		Password: password,
		Database: database,
	}, nil
}

func ProvideServerConfig() (ServerConfig, error) {
	if err := load(); err != nil {
		return ServerConfig{}, err
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return ServerConfig{}, errors.New("SERVER_PORT is not set")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		Host: os.Getenv("SERVER_HOST"),
		Port: uint16(portInt),
	}, nil
}
