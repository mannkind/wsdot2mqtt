package main

import (
	"fmt"
	"strings"

	"github.com/mannkind/twomqtt"
)

type mqttClient struct {
	*twomqtt.MQTTProxy
	mqttClientConfig
	stateUpdateChan stateChannel
}

func newMQTTClient(mqttClientCfg mqttClientConfig, client *twomqtt.MQTTProxy, stateUpdateChan stateChannel) *mqttClient {
	c := mqttClient{
		MQTTProxy:        client,
		mqttClientConfig: mqttClientCfg,
		stateUpdateChan:  stateUpdateChan,
	}

	c.MQTTProxy.
		SetDiscoveryHandler(c.discovery).
		Build()

	return &c
}

func (c *mqttClient) run() {
	c.Run()
	go c.read()
}

func (c *mqttClient) discovery() []twomqtt.MQTTDiscovery {
	mqds := []twomqtt.MQTTDiscovery{}
	if !c.Discovery {
		return mqds
	}

	for _, travelTimeSlug := range c.TravelTimeMapping {
		sensor := strings.ToLower(travelTimeSlug)
		mqd := c.NewMQTTDiscovery("", sensor, "sensor")
		mqd.Icon = "mdi:car"
		mqd.UnitOfMeasurement = "min"
		mqd.Device.Name = Name
		mqd.Device.SWVersion = Version

		mqds = append(mqds, *mqd)
	}

	return mqds
}

func (c *mqttClient) read() {
	for info := range c.stateUpdateChan {
		c.publishState(info)
	}
}

func (c *mqttClient) publishState(info wsdotTravelTime) {
	travelTimeID := fmt.Sprintf("%d", info.TravelTimeID)
	travelTimeSlug := c.TravelTimeMapping[travelTimeID]

	topic := c.StateTopic("", travelTimeSlug)
	payload := fmt.Sprintf("%d", info.CurrentTime)

	if info.Distance == 0 {
		payload = "Closed"
	} else if info.CurrentTime == 0 {
		payload = "Unknown"
	}

	c.Publish(topic, payload)
}
