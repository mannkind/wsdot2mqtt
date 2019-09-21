package main

import (
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
)

type mqttClient struct {
	mqttClientConfig
	*twomqtt.MQTTProxy
	stateUpdateChan stateChannel
}

func newMQTTClient(mqttClientCfg mqttClientConfig, client *twomqtt.MQTTProxy, stateUpdateChan stateChannel) *mqttClient {
	c := mqttClient{
		MQTTProxy:        client,
		mqttClientConfig: mqttClientCfg,
		stateUpdateChan:  stateUpdateChan,
	}

	c.Initialize(
		c.onConnect,
		c.onDisconnect,
	)

	c.LogSettings()

	return &c
}

func (c *mqttClient) run() {
	c.Run()
	go c.receive()
}

func (c *mqttClient) onConnect(client mqtt.Client) {
	log.Info("Connected to MQTT")
	c.Publish(c.AvailabilityTopic(), "online")
	c.publishDiscovery()
}

func (c *mqttClient) onDisconnect(client mqtt.Client, err error) {
	log.WithFields(log.Fields{
		"error": err,
	}).Error("Disconnected from MQTT")
}

func (c *mqttClient) publishDiscovery() {
	if !c.Discovery {
		return
	}

	log.Info("MQTT discovery publishing")

	for _, travelTimeSlug := range c.TravelTimeMapping {
		sensor := strings.ToLower(travelTimeSlug)
		mqd := c.NewMQTTDiscovery("", sensor, "sensor")
		mqd.Icon = "mdi:car"
		mqd.UnitOfMeasurement = "min"

		c.PublishDiscovery(mqd)
	}

	log.Info("Finished MQTT discovery publishing")
}

func (c *mqttClient) receive() {
	for info := range c.stateUpdateChan {
		c.receiveState(info)
	}
}

func (c *mqttClient) receiveState(info wsdotTravelTime) {
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
