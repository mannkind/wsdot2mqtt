package main

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

const defaultDiscoveryName = "wsdot"
const defaultTopicPrefix = "home/wsdot"
const knownTravelID = "132"
const knownTravelIDName = "seattle2everett"
const knownDiscoveryName = "wsdotDiscoveryName"
const knownTopicPrefix = "home/wsdotMQTTTopicPrefix"

func init() {
	log.SetLevel(log.PanicLevel)
}

func setEnvs(d, dn, tp, a string) {
	os.Setenv("MQTT_DISCOVERY", d)
	os.Setenv("MQTT_DISCOVERYNAME", dn)
	os.Setenv("MQTT_TOPICPREFIX", tp)
	os.Setenv("WSDOT_SECRET", "")
	os.Setenv("WSDOT_TRAVELTIMEMAPPING", a)
}

func clearEnvs() {
	setEnvs("false", "", "", "")
}

func TestDiscovery(t *testing.T) {
	defer clearEnvs()

	var tests = []struct {
		TravelTimes     string
		DiscoveryName   string
		TopicPrefix     string
		ExpectedTopic   string
		ExpectedPayload string
	}{
		{
			knownTravelID + ":" + knownTravelIDName,
			defaultDiscoveryName,
			defaultTopicPrefix,
			"homeassistant/sensor/" + defaultDiscoveryName + "/" + knownTravelIDName + "/config",
			"{\"availability_topic\":\"" + defaultTopicPrefix + "/status\",\"device\":{\"identifiers\":[\"" + defaultTopicPrefix + "/status\"],\"manufacturer\":\"twomqtt\",\"name\":\"x2mqtt\",\"sw_version\":\"X.X.X\"},\"icon\":\"mdi:car\",\"name\":\"" + defaultDiscoveryName + " " + knownTravelIDName + "\",\"state_topic\":\"" + defaultTopicPrefix + "/" + knownTravelIDName + "/state\",\"unique_id\":\"" + defaultDiscoveryName + "." + knownTravelIDName + "\",\"unit_of_measurement\":\"min\"}",
		},
		{
			knownTravelID + ":" + knownTravelIDName,
			knownDiscoveryName,
			knownTopicPrefix,
			"homeassistant/sensor/" + knownDiscoveryName + "/" + knownTravelIDName + "/config",
			"{\"availability_topic\":\"" + knownTopicPrefix + "/status\",\"device\":{\"identifiers\":[\"" + knownTopicPrefix + "/status\"],\"manufacturer\":\"twomqtt\",\"name\":\"x2mqtt\",\"sw_version\":\"X.X.X\"},\"icon\":\"mdi:car\",\"name\":\"" + knownDiscoveryName + " " + knownTravelIDName + "\",\"state_topic\":\"" + knownTopicPrefix + "/" + knownTravelIDName + "/state\",\"unique_id\":\"" + knownDiscoveryName + "." + knownTravelIDName + "\",\"unit_of_measurement\":\"min\"}",
		},
	}

	for _, v := range tests {
		setEnvs("true", v.DiscoveryName, v.TopicPrefix, v.TravelTimes)

		c := initialize()
		c.mqttClient.publishDiscovery()

		actualPayload := c.mqttClient.LastPublishedOnTopic(v.ExpectedTopic)
		if actualPayload != v.ExpectedPayload {
			t.Errorf("Actual:%s\nExpected:%s", actualPayload, v.ExpectedPayload)
		}
	}

	os.Setenv("WSDOT_TRAVELTIMEMAPPING", "")
}

func TestReceieveState(t *testing.T) {
	defer clearEnvs()

	var tests = []struct {
		TravelTimes     string
		TravelTimeID    int
		Time            int
		TopicPrefix     string
		ExpectedTopic   string
		ExpectedPayload string
	}{
		{
			knownTravelID + ":" + knownTravelIDName,
			132,
			20,
			defaultTopicPrefix,
			defaultTopicPrefix + "/" + knownTravelIDName + "/state",
			"20",
		},
	}

	for _, v := range tests {
		setEnvs("false", "", v.TopicPrefix, v.TravelTimes)

		obj := wsdotTravelTime{
			TravelTimeID: v.TravelTimeID,
			CurrentTime:  v.Time,
			Distance:     1,
		}

		c := initialize()
		c.mqttClient.receiveState(obj)

		actualPayload := c.mqttClient.LastPublishedOnTopic(v.ExpectedTopic)
		if actualPayload != v.ExpectedPayload {
			t.Errorf("Actual:%s\nExpected:%s", actualPayload, v.ExpectedPayload)
		}
	}
}
