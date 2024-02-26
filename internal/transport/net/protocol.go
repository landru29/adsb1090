package net

import (
	"fmt"
	"strings"
)

type protocolType string

const (
	protocolTypeTCP protocolType = "tcp"
	protocolTypeUDP protocolType = "udp"
)

type protocolDirection string

const (
	protocolDial protocolDirection = "dial"
	protocolBind protocolDirection = "bind"

	defaultProtocolFormat = "nmea"

	defaultAddr = "0.0.0.0:30003"
)

// ProtocolConfig is the net protocol parameter.
type ProtocolConfig struct {
	Addr         string
	Format       string
	Direction    protocolDirection
	ProtocolType protocolType
}

// NewProtocol creates a new ProtocolConfig.
func NewProtocol(pType string) ProtocolConfig {
	return ProtocolConfig{
		ProtocolType: protocolType(pType),
	}
}

// String implements the pflag.Value interface.
func (p *ProtocolConfig) String() string {
	return fmt.Sprintf(
		"%s/%s:%s@%s",
		p.Direction,
		p.ProtocolType,
		p.Format,
		p.Addr,
	)
}

// Set implements the pflag.Value interface.
func (p *ProtocolConfig) Set(str string) error {
	actionSplitter := strings.Split(str, ">")
	switch len(actionSplitter) {
	case 1:
		format, addr := parseData(actionSplitter[0])

		p.Format = format
		p.Direction = protocolDial
		p.Addr = addr

		return nil
	case 2: //nolint: gomnd
		format, addr := parseData(actionSplitter[1])

		p.Format = format
		p.Direction = protocolDirection(actionSplitter[0])
		p.Addr = addr

		return nil
	}

	return fmt.Errorf("wrong format %s (should be like dial>text@0.0.0.0:30003)", str)
}

// Type implements the pflag.Value interface.
func (p *ProtocolConfig) Type() string {
	return "protocol configuration"
}

// IsValid checks if the protocol configuration is valid.
func (p ProtocolConfig) IsValid() bool {
	return p.Addr != ""
}

func parseData(str string) (string, string) {
	addr := defaultAddr
	format := defaultProtocolFormat

	if str != "" {
		addr = str
	}

	splitter := strings.Split(str, "@")
	if len(splitter) > 1 {
		format = splitter[0]
		addr = strings.Join(splitter[1:], "@")
	}

	return format, addr
}
