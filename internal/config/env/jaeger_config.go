package env

import (
	"errors"
	"os"

	"github.com/biryanim/wb_tech_L0/internal/config"
)

const (
	jaegerURLEnvName     = "JAEGER_URL"
	jaegerServiceEnvName = "JAEGER_SERVICE"
)

type jaegerConfig struct {
	url     string
	service string
}

// NewJaegerConfig creates a new JaegerConfig instance by reading JAEGER_URL and JAEGER_SERVICE environment variables.
func NewJaegerConfig() (config.JaegerConfig, error) {
	url := os.Getenv(jaegerURLEnvName)
	if len(url) == 0 {
		return nil, errors.New("jaeger url not found")
	}

	serviceName := os.Getenv(jaegerServiceEnvName)
	if len(serviceName) == 0 {
		return nil, errors.New("jaeger service name not found")
	}

	return &jaegerConfig{
		url:     url,
		service: serviceName,
	}, nil
}

// URL returns the Jaeger collector endpoint URL.
func (j *jaegerConfig) URL() string {
	return j.url
}

// ServiceName returns the service name for Jaeger tracing.
func (j *jaegerConfig) ServiceName() string {
	return j.service
}
