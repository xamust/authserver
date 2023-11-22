package config

import "os"

func Env() (string, bool) {
	confFile := os.Getenv(confEnvFile)
	if confFile == "" {
		return "", false
	}
	return confFile, true
}
