//+build wireinject

package main

import (
	"github.com/google/wire"
	mqttExtCfg "github.com/mannkind/paho.mqtt.golang.ext/cfg"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
)

// Initialize - Compile-time DI
func Initialize() *Wsdot2Mqtt {
	wire.Build(mqttExtCfg.NewMQTTConfig, NewConfig, mqttExtDI.NewMQTTFuncWrapper, NewWsdot2Mqtt)

	return &Wsdot2Mqtt{}
}
