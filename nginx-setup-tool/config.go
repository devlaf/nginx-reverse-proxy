package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	LetsEncryptSettings struct {
		CertName string `json:"cert_name"`
		Email    string `json:"email"`
	} `json:"letsencrypt_settings"`
	Sites   []SiteConfig   `json:"sites"`
	Streams []StreamConfig `json:"streams"`
}

type SiteConfig struct {
	Domain                       string `json:"domain"`
	ProxyToAddress               string `json:"proxy_to_address"`
	ProxyTargetIsDockerContainer bool   `json:"proxy_target_is_docker_container"`
	ForwardHttpToHttps           bool   `json:"forward_http_to_https"`
	IncludeAndRedirectWWW        bool   `json:"include_and_redirect_www"`
	RestrictAccessToIP           string `json:"restrict_access_to_ip"`
}

type StreamConfig struct {
	HostPort                     int    `json:"host_port"`
	ProxyToAddress               string `json:"proxy_to_address"`
	ProxyToPort                  int    `json:"proxy_to_port"`
	ProxyTargetIsDockerContainer bool   `json:"proxy_target_is_docker_container"`
	RestrictAccessToIP           string `json:"restrict_access_to_ip"`
}

func loadConfig(file string) (*Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getDomainList(config *Config) string {
	domains := make([]string, 0)

	for _, site := range config.Sites {
		domains = append(domains, site.Domain)
		if site.IncludeAndRedirectWWW {
			domains = append(domains, fmt.Sprintf("www.%v", site.Domain))
		}
	}

	return strings.Join(domains, ",")
}
