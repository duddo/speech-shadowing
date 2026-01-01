package main

import "os"

type Configuration struct {
	Host       string
	Port       string
	DataPath   string
	ConfigPath string
	TmpPath    string
}

func LoadConfig() *Configuration {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	rootPath := "./" // TODO per ora sempre in dev
	isDev := os.Getenv("DEVENV")
	if isDev == "true" {
		rootPath = "./"
	}

	return &Configuration{
		Host:       host,
		Port:       port,
		DataPath:   rootPath + "data",
		ConfigPath: rootPath + "config",
		TmpPath:    rootPath + "tmp",
	}
}

func (s *Configuration) Addr() string {
	return s.Host + ":" + s.Port
}
