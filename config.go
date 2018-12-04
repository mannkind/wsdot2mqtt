package main

import (
	"log"
	"time"

	"github.com/caarlos0/env"
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
)

// Config - Structured configuration for the application.
type Config struct {
	MQTT              *mqttExtCfg.MQTTConfig
	Secret            string        `env:"WSDOT_SECRET,required"`
	LookupInterval    time.Duration `env:"WSDOT_LOOKUPINTERVAL"    envDefault:"3m"`
	TravelTimeMapping []string      `env:"WSDOT_TRAVELTIMEMAPPING" envDefault:"132:seattle2everett,31:seattle2renton"`
}

// NewConfig - Returns a new reference to a fully configured object.
func NewConfig(mqttCfg *mqttExtCfg.MQTTConfig) *Config {
	c := Config{}
	c.MQTT = mqttCfg

	if c.MQTT.ClientID == "" {
		c.MQTT.ClientID = "DefaultWSDOT2MQTTClientID"
	}

	if c.MQTT.DiscoveryName == "" {
		c.MQTT.DiscoveryName = "wsdot"
	}

	if c.MQTT.TopicPrefix == "" {
		c.MQTT.TopicPrefix = "home/wsdot"
	}

	if err := env.Parse(&c); err != nil {
		log.Printf("Error unmarshaling configuration: %s", err)
	}

	redactedPassword := ""
	if len(c.MQTT.Password) > 0 {
		redactedPassword = "<REDACTED>"
	}

	log.Printf("Environmental Settings:")
	log.Printf("  * ClientID          : %s", c.MQTT.ClientID)
	log.Printf("  * Broker            : %s", c.MQTT.Broker)
	log.Printf("  * Username          : %s", c.MQTT.Username)
	log.Printf("  * Password          : %s", redactedPassword)
	log.Printf("  * Discovery         : %t", c.MQTT.Discovery)
	log.Printf("  * DiscoveryPrefix   : %s", c.MQTT.DiscoveryPrefix)
	log.Printf("  * DiscoveryName     : %s", c.MQTT.DiscoveryName)
	log.Printf("  * TopicPrefix       : %s", c.MQTT.TopicPrefix)
	log.Printf("  * LookupInterval    : %s", c.LookupInterval)
	log.Printf("  * TravelTimeMapping : %s", c.TravelTimeMapping)

	return &c
}
