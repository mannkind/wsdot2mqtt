package main

import (
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
	mqttExtHA "github.com/mannkind/paho.mqtt.golang.ext/ha"
	log "github.com/sirupsen/logrus"
)

const (
	sensorUniqueTemplate = "%s.%s"
	sensorTopicTemplate  = "%s/%s/state"
)

// mqttClient - Lookup travel time information on wsdot.wa.gov.
type mqttClient struct {
	discovery       bool
	discoveryPrefix string
	discoveryName   string
	topicPrefix     string
	travelTimes     map[string]string

	client        mqtt.Client
	lastPublished map[string]string
}

func newMQTTClient(config *config, mqttFuncWrapper *mqttExtDI.MQTTFuncWrapper) *mqttClient {
	c := mqttClient{
		discovery:       config.MQTT.Discovery,
		discoveryPrefix: config.MQTT.DiscoveryPrefix,
		discoveryName:   config.MQTT.DiscoveryName,
		topicPrefix:     config.MQTT.TopicPrefix,

		travelTimes:   map[string]string{},
		lastPublished: map[string]string{},
	}

	// Create a mapping between travel time ids and names
	for _, m := range config.TravelTimeMapping {
		parts := strings.Split(m, ":")
		if len(parts) != 2 {
			continue
		}

		travelTimeID := parts[0]
		travelTimeName := parts[1]
		c.travelTimes[travelTimeID] = travelTimeName
	}

	opts := mqttFuncWrapper.
		ClientOptsFunc().
		AddBroker(config.MQTT.Broker).
		SetClientID(config.MQTT.ClientID).
		SetOnConnectHandler(c.onConnect).
		SetConnectionLostHandler(c.onDisconnect).
		SetUsername(config.MQTT.Username).
		SetPassword(config.MQTT.Password).
		SetWill(c.availabilityTopic(), "offline", 0, true)

	c.client = mqttFuncWrapper.ClientFunc(opts)

	return &c
}

func (c *mqttClient) run() {
	c.runAfter(0 * time.Second)
}

func (c *mqttClient) runAfter(delay time.Duration) {
	time.Sleep(delay)

	log.Info("Connecting to MQTT")
	if token := c.client.Connect(); !token.Wait() || token.Error() != nil {
		log.WithFields(log.Fields{
			"error": token.Error(),
		}).Error("Error connecting to MQTT")

		delay = c.adjustReconnectDelay(delay)

		log.WithFields(log.Fields{
			"delay": delay,
		}).Info("Sleeping before attempting to reconnect to MQTT")

		c.runAfter(delay)
	}
}

func (c *mqttClient) adjustReconnectDelay(delay time.Duration) time.Duration {
	var maxDelay float64 = 120
	defaultDelay := 2 * time.Second

	// No delay, set to default delay
	if delay.Seconds() == 0 {
		delay = defaultDelay
	} else {
		// Increment the delay
		delay = delay * 2

		// If the delay is above two minutes, reset to default
		if delay.Seconds() > maxDelay {
			delay = defaultDelay
		}
	}

	return delay
}

func (c *mqttClient) onConnect(client mqtt.Client) {
	log.Info("Connected to MQTT")
	c.publish(c.availabilityTopic(), "online")
	c.publishDiscovery()
}

func (c *mqttClient) onDisconnect(client mqtt.Client, err error) {
	log.WithFields(log.Fields{
		"error": err,
	}).Error("Disconnected from MQTT")
}

func (c *mqttClient) availabilityTopic() string {
	return fmt.Sprintf("%s/status", c.topicPrefix)
}

func (c *mqttClient) publishDiscovery() {
	if !c.discovery {
		return
	}

	for _, travelTimeSlug := range c.travelTimes {
		sensor := strings.ToLower(travelTimeSlug)
		mqd := mqttExtHA.MQTTDiscovery{
			DiscoveryPrefix: c.discoveryPrefix,
			Component:       "sensor",
			NodeID:          c.discoveryName,
			ObjectID:        sensor,

			AvailabilityTopic: c.availabilityTopic(),
			Name:              fmt.Sprintf("%s %s", c.discoveryName, sensor),
			StateTopic:        fmt.Sprintf(sensorTopicTemplate, c.topicPrefix, sensor),
			UniqueID:          fmt.Sprintf("%s.%s", c.discoveryName, sensor),
			Icon:              "mdi:car",
			UnitOfMeasurement: "min",
		}

		mqd.PublishDiscovery(c.client)
	}
}

func (c *mqttClient) receiveCommand(int64, event) {}
func (c *mqttClient) receiveState(e event) {
	travelTimeSlug := e.key
	response := e.data

	topic := fmt.Sprintf(sensorTopicTemplate, c.topicPrefix, travelTimeSlug)
	payload := fmt.Sprintf("%d", response.CurrentTime)

	if response.Distance == 0 {
		payload = "Closed"
	} else if response.CurrentTime == 0 {
		payload = "Unknown"
	}

	c.publish(topic, payload)
}

func (c *mqttClient) publish(topic string, payload string) {
	llog := log.WithFields(log.Fields{
		"topic":   topic,
		"payload": payload,
	})
	// Should we publish this again?
	// NOTE: We must allow the availability topic to publish duplicates
	if lastPayload, ok := c.lastPublished[topic]; topic != c.availabilityTopic() && ok && lastPayload == payload {
		llog.Debug("Duplicate payload")
		return
	}

	llog.Info("Publishing to MQTT")

	retain := true
	if token := c.client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Error("Publishing error")
	}

	llog.Debug("Published to MQTT")
	c.lastPublished[topic] = payload
}
