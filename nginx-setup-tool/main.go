package main

//go:generate go run nginx-templates/include-templates.go -templatedir=./nginx-templates

import (
	"fmt"
	"os"
)

func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func succeed(msg string) {
	fmt.Fprintln(os.Stdout, msg)
	os.Exit(0)
}

func main() {
	args := fetchArgs()

	config, err := loadConfig(args.ConfigFilepath)
	if err != nil {
		die(fmt.Sprintf("ERR -- could not read config: %v", err))
	}

	switch args.Operation {
	case "certname":
		succeed(config.LetsEncryptSettings.CertName)
	case "email":
		succeed(config.LetsEncryptSettings.Email)
	case "domains":
		succeed(getDomainList(config))
	case "write":
		translateToNginxConfig(config)
	}
}
