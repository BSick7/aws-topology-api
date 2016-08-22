package types

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
)

type AgentConfig struct {
	Bind AgentConfigBind `hcl:"bind"`
}

func NewAgentConfigFromLocation(location string) (*AgentConfig, error) {
	if location == "" {
		return &AgentConfig{}, nil
	}

	if s, err := os.Stat(location); err != nil {
		return nil, err
	} else if s.IsDir() {
		files, err := filepath.Glob(filepath.Join(location, "*.hcl"))
		if err != nil {
			return nil, err
		}

		ac := &AgentConfig{}
		for _, file := range files {
			if acn, err := NewAgentConfigFromFile(file); err != nil {
				return nil, err
			} else {
				ac.Merge(acn)
			}
		}
		return ac, nil
	} else {
		return NewAgentConfigFromFile(location)
	}
}

func NewAgentConfigFromFile(filename string) (*AgentConfig, error) {
	ac := &AgentConfig{}
	if in, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		err := hcl.Decode(ac, string(in))
		return ac, err
	}
}

func (ac *AgentConfig) Merge(acn *AgentConfig) {
	if acn.Bind.Address != "" {
		ac.Bind.Address = acn.Bind.Address
	}
	if acn.Bind.Port != 0 {
		ac.Bind.Port = acn.Bind.Port
	}
}

func (ac *AgentConfig) Validate() error {
	if ac.Bind.Address == "0.0.0.0" {
		if eip, err := externalIP(); err != nil {
			return err
		} else {
			ac.Bind.Address = eip
		}
	}

	return nil
}

type AgentConfigBind struct {
	Address string `hcl:"address"`
	Port    int    `hcl:"port"`
}

func (b AgentConfigBind) String() string {
	return fmt.Sprintf("%s:%d", b.Address, b.Port)
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("could not find external interface")
}
