package main

import (
	"time"

	"github.com/caarlos0/env"
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	log "github.com/sirupsen/logrus"
)

type config struct {
	MQTT              *mqttExtCfg.MQTTConfig
	Secret            string        `env:"WSDOT_SECRET,required"`
	LookupInterval    time.Duration `env:"WSDOT_LOOKUPINTERVAL"    envDefault:"3m"`
	TravelTimeMapping []string      `env:"WSDOT_TRAVELTIMEMAPPING" envDefault:"132:seattle2everett,31:seattle2renton"`
	DebugLogLevel     bool          `env:"WSDOT_DEBUG" envDefault:"false"`
}

func newConfig(mqttCfg *mqttExtCfg.MQTTConfig) *config {
	c := config{}
	c.MQTT = mqttCfg
	c.MQTT.Defaults("DefaultWSDOT2MQTTClientID", "wsdot", "home/wsdot")

	if err := env.Parse(&c); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to unmarshal configuration")
	}

	log.WithFields(log.Fields{
		"WSDOT.LookupInterval": c.LookupInterval,
		"WSDOT.TravelMapping":  c.TravelTimeMapping,
		"WSDOT.DebugLogLevel":  c.DebugLogLevel,
	}).Info("Environmental Settings")

	if c.DebugLogLevel {
		log.SetLevel(log.DebugLevel)
		log.Debug("Enabling the debug log level")
	}

	return &c
}
