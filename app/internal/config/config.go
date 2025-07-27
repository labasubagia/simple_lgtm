package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppName                  string
	Port                     int
	OtelExporterOLTPEndpoint string
}

func Load() *Config {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "app"
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 5000 // Default port
	}

	otelExporterOLTPEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otelExporterOLTPEndpoint == "" {
		otelExporterOLTPEndpoint = "http://localhost:4318"
	}

	return &Config{
		AppName:                  appName,
		Port:                     port,
		OtelExporterOLTPEndpoint: otelExporterOLTPEndpoint,
	}
}
