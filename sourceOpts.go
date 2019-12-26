package main

import (
	"time"
)

type sourceOpts struct {
	globalOpts
	Secret         string        `env:"WSDOT_SECRET,required"`
	LookupInterval time.Duration `env:"WSDOT_LOOKUPINTERVAL"    envDefault:"3m"`
}
