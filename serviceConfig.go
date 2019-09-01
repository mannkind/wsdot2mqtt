package main

import (
	"time"
)

type serviceClientConfig struct {
	globalClientConfig
	Secret         string        `env:"WSDOT_SECRET,required"`
	LookupInterval time.Duration `env:"WSDOT_LOOKUPINTERVAL"    envDefault:"3m"`
}
