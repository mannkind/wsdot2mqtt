package mqtt

import (
	"fmt"
	"strings"

	"github.com/mannkind/twomqtt"
	"github.com/mannkind/wsdot2mqtt/shared"
)

// Writer is for writing a shared representation to MQTT
type Writer struct {
	*twomqtt.MQTT
	opts     Opts
	incoming <-chan shared.Representation
}

// NewWriter creates a new Writer for writing a shared representation to MQTT
func NewWriter(mqtt *twomqtt.MQTT, opts Opts, incoming <-chan shared.Representation) *Writer {
	c := Writer{
		MQTT:     mqtt,
		opts:     opts,
		incoming: incoming,
	}

	c.MQTT.
		SetDiscoveryHandler(c.discovery).
		SetReadIncomingChannelHandler(c.read).
		Initialize()

	return &c
}

// discovery objects to publish when MQTT discovery is enabled
func (c *Writer) discovery() []twomqtt.MQTTDiscovery {
	mqds := []twomqtt.MQTTDiscovery{}
	if !c.Discovery {
		return mqds
	}

	for _, travelTimeSlug := range c.opts.TravelTimeMapping {
		deviceName := ""
		sensorName := strings.ToLower(travelTimeSlug)
		sensorType := "sensor"

		mqd := twomqtt.NewMQTTDiscovery(c.opts.MQTTOpts, deviceName, sensorName, sensorType)
		mqd.Icon = "mdi:car"
		mqd.UnitOfMeasurement = "min"
		mqd.Device.Name = shared.Name
		mqd.Device.SWVersion = shared.Version

		mqds = append(mqds, *mqd)
	}

	return mqds
}

// read incoming shared representations and publish them to MQTT
func (c *Writer) read() {
	for info := range c.incoming {
		c.publish(info)
	}
}

// publish a shared representation to MQTT
func (c *Writer) publish(info shared.Representation) twomqtt.MQTTMessage {
	travelTimeID := fmt.Sprintf("%d", info.TravelTimeID)
	travelTimeSlug := c.opts.TravelTimeMapping[travelTimeID]

	topic := c.StateTopic("", travelTimeSlug)
	payload := fmt.Sprintf("%d", info.CurrentTime)

	if info.Distance == 0 {
		payload = "Closed"
	} else if info.CurrentTime == 0 {
		payload = "Unknown"
	}

	return c.Publish(topic, payload)
}
