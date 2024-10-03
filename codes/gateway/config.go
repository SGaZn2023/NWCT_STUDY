package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenerFile string `yaml:"listener_file"`
}

func ParseConfig(confFile string) (*Config, error) {
	content, err := os.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type ListenerConfig struct {
	ClientID         string `json:"client_id"`
	PublicIP         string `json:"public_ip"`
	PublicPort       uint16 `json:"public_port"`
	PublicProtocol   string `json:"public_protocol"`
	InternalIP       string `json:"internal_ip"`
	InternalPort     uint16 `json:"internal_port"`
	InternalProtocol string `json:"internal_protocol"`
}

func ParseListenerConfig(confFile string) ([]*ListenerConfig, error) {
	content, err := os.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	var cfg = make([]*ListenerConfig, 0)

	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
