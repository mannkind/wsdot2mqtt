package main

import (
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
	"testing"
)

func defaultWsdot2Mqtt() *Wsdot2Mqtt {
	c := NewWsdot2Mqtt(NewConfig(mqttExtCfg.NewMQTTConfig()), mqttExtDI.NewMQTTFuncWrapper())
	return c
}

func TestMqttConnect(t *testing.T) {
	c := defaultWsdot2Mqtt()
	c.onConnect(c.client)
}
