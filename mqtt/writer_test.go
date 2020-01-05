package mqtt

import (
	"os"
	"testing"

	"github.com/mannkind/twomqtt"
	"github.com/mannkind/wsdot2mqtt/shared"
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

func initialize() *Writer {
	opts := shared.NewOpts()
	v := shared.NewRepresentationChannel()
	v3 := shared.NewRepresentationChannelIncoming(v)
	mqttOpts := NewOpts(opts)
	twomqttMQTTOpts := mqttOpts.MQTTOpts
	twomqttMQTT := twomqtt.NewMQTT(twomqttMQTTOpts)
	writer := NewWriter(twomqttMQTT, mqttOpts, v3)
	return writer
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
		TravelTimes           string
		DiscoveryName         string
		TopicPrefix           string
		ExpectedName          string
		ExpectedStateTopic    string
		ExpectedUniqueID      string
		ExpectedIcon          string
		ExpectedUnitOfMeasure string
	}{
		{
			knownTravelID + ":" + knownTravelIDName,
			defaultDiscoveryName,
			defaultTopicPrefix,
			defaultDiscoveryName + " " + knownTravelIDName,
			defaultTopicPrefix + "/" + knownTravelIDName + "/state",
			defaultDiscoveryName + "." + knownTravelIDName,
			"mdi:car",
			"min",
		},
		{
			knownTravelID + ":" + knownTravelIDName,
			knownDiscoveryName,
			knownTopicPrefix,
			knownDiscoveryName + " " + knownTravelIDName,
			knownTopicPrefix + "/" + knownTravelIDName + "/state",
			knownDiscoveryName + "." + knownTravelIDName,
			"mdi:car",
			"min",
		},
	}

	for _, v := range tests {
		setEnvs("true", v.DiscoveryName, v.TopicPrefix, v.TravelTimes)

		c := initialize()
		mqds := c.discovery()

		for _, mqd := range mqds {
			if mqd.Name != v.ExpectedName {
				t.Errorf("discovery Name does not match; %s vs %s", mqd.Name, v.ExpectedName)
			}
			if mqd.StateTopic != v.ExpectedStateTopic {
				t.Errorf("discovery StateTopic does not match; %s vs %s", mqd.StateTopic, v.ExpectedStateTopic)
			}
			if mqd.UniqueID != v.ExpectedUniqueID {
				t.Errorf("discovery UniqueID does not match; %s vs %s", mqd.UniqueID, v.ExpectedUniqueID)
			}
			if mqd.Icon != v.ExpectedIcon {
				t.Errorf("discovery Icon does not match; %s vs %s", mqd.Icon, v.ExpectedIcon)
			}
			if mqd.UnitOfMeasurement != v.ExpectedUnitOfMeasure {
				t.Errorf("discovery UnitOfMeasurement does not match; %s vs %s", mqd.UnitOfMeasurement, v.ExpectedUnitOfMeasure)
			}
		}
	}

	os.Setenv("WSDOT_TRAVELTIMEMAPPING", "")
}

func TestPublish(t *testing.T) {
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

		obj := shared.Representation{
			TravelTimeID: v.TravelTimeID,
			CurrentTime:  v.Time,
			Distance:     1,
		}

		c := initialize()
		published := c.publish(obj)

		if published.Payload != v.ExpectedPayload {
			t.Errorf("Actual:%s\nExpected:%s", published.Payload, v.ExpectedPayload)
		}
	}
}
