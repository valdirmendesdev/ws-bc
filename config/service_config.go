package config

import "fmt"

type ServiceConfig struct {
	host string
	port string
}

func New(host, port string) *ServiceConfig {
	return &ServiceConfig{
		host: host,
		port: port,
	}
}

func (c *ServiceConfig) Host() string {
	if c.host == "" {
		return "localhost"
	}
	return c.host
}

func (c *ServiceConfig) Port() string {
	if c.port == "" {
		return "8080"
	}
	return c.port
}

func (c *ServiceConfig) FullHost() string {
	return fmt.Sprintf("%s:%s", c.Host(), c.Port())
}
