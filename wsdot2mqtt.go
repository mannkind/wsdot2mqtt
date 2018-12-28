package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqttExtDI "github.com/mannkind/paho.mqtt.golang.ext/di"
	mqttExtHA "github.com/mannkind/paho.mqtt.golang.ext/ha"
	resty "gopkg.in/resty.v1"
)

const (
	travelTimeURL       = "http://www.wsdot.wa.gov/Traffic/api/TravelTimes/TravelTimesREST.svc/GetTravelTimeAsJson"
	sensorTopicTemplate = "%s/%s/state"
)

// Wsdot2Mqtt - Lookup travel time information on wsdot.wa.gov.
type Wsdot2Mqtt struct {
	discovery       bool
	discoveryPrefix string
	discoveryName   string
	topicPrefix     string
	lookupInterval  time.Duration
	travelTimes     map[string]string
	secret          string

	client mqtt.Client
}

// NewWsdot2Mqtt - Returns a new reference to a fully configured object.
func NewWsdot2Mqtt(config *Config, mqttFuncWrapper *mqttExtDI.MQTTFuncWrapper) *Wsdot2Mqtt {
	x := Wsdot2Mqtt{
		discovery:       config.MQTT.Discovery,
		discoveryPrefix: config.MQTT.DiscoveryPrefix,
		discoveryName:   config.MQTT.DiscoveryName,
		topicPrefix:     config.MQTT.TopicPrefix,
		lookupInterval:  config.LookupInterval,
		secret:          config.Secret,
	}
	x.travelTimes = make(map[string]string, 0)

	// Create a mapping between travel time ids and names
	for _, m := range config.TravelTimeMapping {
		parts := strings.Split(m, ":")
		if len(parts) != 2 {
			continue
		}

		travelTimeID := parts[0]
		travelTimeName := parts[1]
		x.travelTimes[travelTimeID] = travelTimeName
	}

	opts := mqttFuncWrapper.
		ClientOptsFunc().
		AddBroker(config.MQTT.Broker).
		SetClientID(config.MQTT.ClientID).
		SetOnConnectHandler(x.onConnect).
		SetConnectionLostHandler(x.onDisconnect).
		SetUsername(config.MQTT.Username).
		SetPassword(config.MQTT.Password)

	x.client = mqttFuncWrapper.ClientFunc(opts)

	return &x
}

// Run - Start the collection lookup process
func (t *Wsdot2Mqtt) Run() error {
	log.Print("Connecting to MQTT")
	if token := t.client.Connect(); !token.Wait() || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (t *Wsdot2Mqtt) onConnect(client mqtt.Client) {
	log.Print("Connected to MQTT")

	if !client.IsConnected() {
		log.Print("Subscribe Error: Not Connected (Reloading Config?)")
		return
	}

	if t.discovery {
		t.publishDiscovery()
	}

	go t.loop()
}

func (t *Wsdot2Mqtt) onDisconnect(client mqtt.Client, err error) {
	log.Printf("Disconnected from MQTT: %s.", err)
}

func (t *Wsdot2Mqtt) loop() {
	for {
		for travelTimeID, travelTimeSlug := range t.travelTimes {
			tt, err := t.lookupTravelTime(travelTimeID)
			if err != nil {
				log.Printf("Unable to lookup travel time for TravelTimeID: %s", travelTimeID)
				continue
			}
			t.publishTravelTime(tt, travelTimeSlug)
		}

		time.Sleep(t.lookupInterval)
	}
}

func (t *Wsdot2Mqtt) lookupTravelTime(travelTimeID string) (*travelTimeAPIResponse, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&travelTimeAPIResponse{}).
		SetQueryParams(map[string]string{
			"AccessCode":   t.secret,
			"TravelTimeID": travelTimeID,
		}).
		Get(travelTimeURL)

	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("Unble to lookup the travel time for %s", travelTimeID)
	}

	return resp.Result().(*travelTimeAPIResponse), nil
}

func (t *Wsdot2Mqtt) publishTravelTime(response *travelTimeAPIResponse, travelTimeSlug string) {
	if response.CurrentTime == 0 {
		log.Printf("Ignoring travel time of 0 for %s", travelTimeSlug)
		return
	}

	topic := fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, travelTimeSlug)
	payload := fmt.Sprintf("%d", response.CurrentTime)

	t.publish(topic, payload)
}

func (t *Wsdot2Mqtt) publishDiscovery() {
	for _, travelTimeSlug := range t.travelTimes {
		sensor := strings.ToLower(travelTimeSlug)
		mqd := mqttExtHA.MQTTDiscovery{
			DiscoveryPrefix: t.discoveryPrefix,
			Component:       "sensor",
			NodeID:          t.discoveryName,
			ObjectID:        sensor,

			Name:              fmt.Sprintf("%s %s", t.discoveryName, sensor),
			StateTopic:        fmt.Sprintf(sensorTopicTemplate, t.topicPrefix, sensor),
			UniqueID:          fmt.Sprintf("%s.%s", t.discoveryName, sensor),
			Icon:              "mdi:car",
			UnitOfMeasurement: "min",
		}

		mqd.PublishDiscovery(t.client)
	}
}

func (t *Wsdot2Mqtt) publish(topic string, payload string) {
	retain := true
	if token := t.client.Publish(topic, 0, retain, payload); token.Wait() && token.Error() != nil {
		log.Printf("Publish Error: %s", token.Error())
	}

	log.Print(fmt.Sprintf("Publishing - Topic: %s ; Payload: %s", topic, payload))
}
