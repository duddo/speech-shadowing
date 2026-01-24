package main

import "os"

type Configuration struct {
	Host       string
	Port       string
	DataPath   string
	AudioPath  string
	ConfigPath string
	TmpPath    string
	StaticPath string
}

func NewConfiguration() *Configuration {
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
		AudioPath:  rootPath + "public/audio",
		ConfigPath: rootPath + "config",
		TmpPath:    rootPath + "tmp",
		StaticPath: rootPath + "public",
	}
}

func (s *Configuration) Addr() string {
	return s.Host + ":" + s.Port
}
