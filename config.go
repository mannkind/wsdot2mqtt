package main

import (
	"reflect"

	"github.com/caarlos0/env/v6"
	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
)

type config struct {
	GeneralConfig       twomqtt.GeneralConfig
	GlobalClientConfig  globalClientConfig
	MQTTClientConfig    mqttClientConfig
	ServiceClientConfig serviceClientConfig
}

func newConfig() config {
	c := config{
		GeneralConfig:       twomqtt.GeneralConfig{},
		GlobalClientConfig:  globalClientConfig{},
		MQTTClientConfig:    mqttClientConfig{},
		ServiceClientConfig: serviceClientConfig{},
	}

	// Manually parse the address:name mapping
	if err := env.ParseWithFuncs(&c, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(travelMapping{}): twomqtt.SimpleKVMapParser(":", ","),
	}); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to unmarshal configuration")
	}

	// Defaults
	if c.MQTTClientConfig.MQTTProxyConfig.DiscoveryName == "" {
		c.MQTTClientConfig.MQTTProxyConfig.DiscoveryName = "wsdot"
	}

	if c.MQTTClientConfig.MQTTProxyConfig.TopicPrefix == "" {
		c.MQTTClientConfig.MQTTProxyConfig.TopicPrefix = "home/wsdot"
	}

	// env.Parse* does not seem to work with embedded structs
	c.MQTTClientConfig.TravelTimeMapping = c.GlobalClientConfig.TravelTimeMapping
	c.ServiceClientConfig.TravelTimeMapping = c.GlobalClientConfig.TravelTimeMapping

	if c.GeneralConfig.DebugLogLevel {
		log.SetLevel(log.DebugLevel)
	}

	return c
}
