package main

import (
	"fmt"
	"os"
	"text/template"
)

func translateToNginxConfig(config *Config) {
	if err := ensureNginxConfigPaths(); err != nil {
		die(fmt.Sprintf("Could not create nginx cfg directories: %v", err))
	}

	if err := writeCertPathsConfig(config); err != nil {
		die(fmt.Sprintf("Could not create cert-paths.conf: %v", err))
	}

	if err := writeHttpsDefaultConfig(); err != nil {
		die(fmt.Sprintf("Could not create https-default.conf: %v", err))
	}

	for _, site := range config.Sites {
		if err := writeSiteConfig(site); err != nil {
			die(fmt.Sprintf("Could not create %v.conf: %v", site.Domain, err))
		}
	}

	for _, stream := range config.Streams {
		if err := writeStreamConfig(stream); err != nil {
			die(fmt.Sprintf("Could not create stream_%v.conf: %v", stream.HostPort, err))
		}
	}
}

func writeCertPathsConfig(config *Config) error {
	tmpl, err := template.New("").Parse(template_cert_paths)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("/etc/nginx/sites-include/cert-paths.conf")
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}

func writeHttpsDefaultConfig() error {
	file, err := os.Create("/etc/nginx/sites-enabled/https-default.conf")
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(template_https_default + "\n")
	return nil
}

func writeSiteConfig(config SiteConfig) error {
	tmpl, err := template.New("").Parse(template_site)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fmt.Sprintf("/etc/nginx/sites-enabled/%v.conf", config.Domain))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}

func writeStreamConfig(config StreamConfig) error {
	tmpl, err := template.New("").Parse(template_stream)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(fmt.Sprintf("/etc/nginx/streams-enabled/%v.conf", config.HostPort))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, config)
}

func ensureNginxConfigPaths() error {
	for _, dir := range []string{"/etc/nginx/sites-enabled", "/etc/nginx/streams-enabled", "/etc/nginx/sites-enabled"} {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir(dir, os.ModeDir)
			} else {
				return err
			}
		}
	}
	return nil
}
