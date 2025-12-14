package config

import "flag"

type Config struct {
	Listen   string
	Upstream string
}

func Load() Config {
	up := flag.String("resolver", "8.8.8.8:53", "")
	flag.Parse()

	return Config{
		Listen:   "127.0.0.1:2053",
		Upstream: *up,
	}
}
