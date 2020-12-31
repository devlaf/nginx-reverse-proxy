package main

import (
	"flag"
	"fmt"
)

type Args struct {
	ConfigFilepath string
	Operation      string
}

func fetchArgs() Args {
	var configFilepath string
	var op string

	flag.StringVar(&configFilepath, "config", "", "path to config file")
	flag.StringVar(&op, "op", "", "operation to run: [certname|email|domains|write]")

	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("nginx-setup-tool -op={operation} -config=/my/config.json")
		fmt.Println("\nParameters:")
		flag.PrintDefaults()
		fmt.Println("\nValid Operations:")
		fmt.Println("certname -- grabs the certificate name for certbot to use from config file")
		fmt.Println("email    -- grabs the email for certbot to use from the config file")
		fmt.Println("domains  -- creates a comma-separated list of all domains required by the config file")
		fmt.Println("write    -- writes nginx config files to /etc/nginx/... based on whats in the config ")
	}

	flag.Parse()

	if configFilepath == "" {
		die("config argument is required, this tool doesn't make much sense without one: see [nginx-setup-tool --help] for details.")
	}

	validOps := []string{"certname", "email", "domains", "write"}
	if !stringInCollection(op, validOps) {
		die("Invalid operation specified: see [nginx-setup-tool --help] for details.")
	}

	return Args{
		ConfigFilepath: configFilepath,
		Operation:      op,
	}
}

func stringInCollection(str string, cln []string) bool {
	for _, element := range cln {
		if str == element {
			return true
		}
	}
	return false
}
