package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         int    `json:"port"`
	GPTAPIURL    string `json:"gptapiurl"`
	PythonAPIURL string `json:"pythonapiurl"`
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() Config {
	cfg := Config{
		Port:         8080,
		GPTAPIURL:    "http://gpt-server:5000",
		PythonAPIURL: "http://python-server:5001",
	}

	if os.Getenv("PORT") != "" {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		cfg.Port = port
	}

	if os.Getenv("GPT_API_URL") != "" {
		cfg.GPTAPIURL = os.Getenv("GPT_API_URL")
	}

	if os.Getenv("PYTHON_API_URL") != "" {
		cfg.PythonAPIURL = os.Getenv("PYTHON_API_URL")
	}

	return cfg
}
