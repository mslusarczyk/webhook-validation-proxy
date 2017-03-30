package config

import "github.com/namsral/flag"

type Params struct {
	Context string
	Port    string
	Target  string
	Secret  string
	Cidr    string
}

func ParseParams() *Params {
	params := new(Params)
	flag.String(flag.DefaultConfigFlagname, "", "Path to config file")
	flag.StringVar(&params.Context, "context", "/github-webhook/", "Context path for proxy for webhooks handling")
	flag.StringVar(&params.Port, "port", "8888", "Port for proxy to listen on")
	flag.StringVar(&params.Target, "target", "http://localhost:8080", "Target address with port")
	flag.StringVar(&params.Secret, "secret", "", "Secret assosiated with GH webhook")
	flag.StringVar(&params.Cidr, "cidr", "192.30.252.0/22", "CIDR of GH servers")
	flag.Parse()

	return params
}
