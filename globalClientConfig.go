package main

type globalClientConfig struct {
	TravelTimeMapping travelMapping `env:"WSDOT_TRAVELTIMEMAPPING" envDefault:"132:seattle2everett,31:seattle2renton"`
}
