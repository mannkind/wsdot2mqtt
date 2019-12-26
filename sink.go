package main

import (
	"fmt"
	"strings"

	"github.com/mannkind/twomqtt"
)

type sink struct {
	*twomqtt.MQTT
	config   sinkOpts
	incoming <-chan sourceRep
}

func newSink(mqtt *twomqtt.MQTT, config sinkOpts, incoming <-chan sourceRep) *sink {
	c := sink{
		MQTT:     mqtt,
		config:   config,
		incoming: incoming,
	}

	c.MQTT.
		SetDiscoveryHandler(c.discovery).
		SetReadIncomingChannelHandler(c.read).
		Initialize()

	return &c
}

func (c *sink) run() {
	c.Run()
}

func (c *sink) discovery() []twomqtt.MQTTDiscovery {
	mqds := []twomqtt.MQTTDiscovery{}
	if !c.Discovery {
		return mqds
	}

	for _, travelTimeSlug := range c.config.TravelTimeMapping {
		deviceName := ""
		sensorName := strings.ToLower(travelTimeSlug)
		sensorType := "sensor"

		mqd := twomqtt.NewMQTTDiscovery(c.config.MQTTOpts, deviceName, sensorName, sensorType)
		mqd.Icon = "mdi:car"
		mqd.UnitOfMeasurement = "min"
		mqd.Device.Name = Name
		mqd.Device.SWVersion = Version

		mqds = append(mqds, *mqd)
	}

	return mqds
}

func (c *sink) read() {
	for info := range c.incoming {
		c.publish(info)
	}
}

func (c *sink) publish(info sourceRep) twomqtt.MQTTMessage {
	travelTimeID := fmt.Sprintf("%d", info.TravelTimeID)
	travelTimeSlug := c.config.TravelTimeMapping[travelTimeID]

	topic := c.StateTopic("", travelTimeSlug)
	payload := fmt.Sprintf("%d", info.CurrentTime)

	if info.Distance == 0 {
		payload = "Closed"
	} else if info.CurrentTime == 0 {
		payload = "Unknown"
	}

	return c.Publish(topic, payload)
}
