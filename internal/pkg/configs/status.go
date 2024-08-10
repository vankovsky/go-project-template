package configs

import "github.com/pkg/errors"

type httpServer struct {
	Address string `yaml:"address"`
	Port    uint   `yaml:"port"`
}

func (h *httpServer) initiate() error {
	if err := h.validate(); err != nil {
		return errors.Wrap(err, "validate")
	}

	return nil
}

func (h *httpServer) validate() error {
	if h.Port == 0 {
		return errors.New("port is empty")
	}

	return nil
}
