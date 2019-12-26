package main

type globalOpts struct {
	TravelTimeMapping sourceMapping `env:"WSDOT_TRAVELTIMEMAPPING" envDefault:"132:seattle2everett,31:seattle2renton"`
}
