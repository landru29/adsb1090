package config

import (
	"fmt"
	"strings"
)

// HTTPConfig is a HTTP configuration.
type HTTPConfig struct {
	Addr    string
	APIPath string
}

const defaulthttpAPIpath = "/api"

// String implements the pflag.Value interface.
func (h *HTTPConfig) String() string {
	return fmt.Sprintf("%s%s", h.Addr, h.APIPath)
}

// Set implements the pflag.Value interface.
func (h *HTTPConfig) Set(str string) error {
	splitter := strings.Split(str, "/")
	if len(splitter) > 1 {
		apiPath := strings.Join(splitter[1:], "/")

		*h = HTTPConfig{
			Addr:    splitter[0],
			APIPath: "/" + apiPath,
		}

		return nil
	}

	*h = HTTPConfig{
		APIPath: defaulthttpAPIpath,
		Addr:    str,
	}

	return nil
}

// Type implements the pflag.Value interface.
func (h *HTTPConfig) Type() string {
	return "http configuration"
}
