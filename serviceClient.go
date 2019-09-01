package main

import (
	"reflect"
	"time"

	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

const (
	travelTimeURL = "https://www.wsdot.wa.gov/Traffic/api/TravelTimes/TravelTimesREST.svc/GetTravelTimeAsJson"
)

type serviceClient struct {
	twomqtt.Publisher
	serviceClientConfig
	observers map[twomqtt.Observer]struct{}
}

func newServiceClient(serviceClientCfg serviceClientConfig) *serviceClient {
	c := serviceClient{
		serviceClientConfig: serviceClientCfg,
		observers:           map[twomqtt.Observer]struct{}{},
	}

	log.WithFields(log.Fields{
		"WSDOT.LookupInterval": c.LookupInterval,
		"WSDOT.TravelMapping":  c.TravelTimeMapping,
	}).Info("Service Environmental Settings")

	return &c
}

func (c *serviceClient) run() {
	go c.loop()
}

func (c *serviceClient) Register(l twomqtt.Observer) {
	c.observers[l] = struct{}{}
}

func (c *serviceClient) sendState(e twomqtt.Event) {
	log.WithFields(log.Fields{
		"event": e,
	}).Debug("Sending event to observers")

	for o := range c.observers {
		o.ReceiveState(e)
	}

	log.Debug("Finished sending event to observers")
}

func (c *serviceClient) loop() {
	for {
		log.Info("Looping")
		for travelTimeID := range c.TravelTimeMapping {
			info, err := c.lookup(travelTimeID)
			if err != nil {
				continue
			}

			event, err := c.adapt(info)
			if err != nil {
				continue
			}

			c.sendState(event)
		}

		log.WithFields(log.Fields{
			"sleep": c.LookupInterval,
		}).Info("Finished looping; sleeping")
		time.Sleep(c.LookupInterval)
	}
}

func (c *serviceClient) lookup(travelTimeID string) (*wsdotTravelTimeAPIResponse, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&wsdotTravelTimeAPIResponse{}).
		SetQueryParams(map[string]string{
			"AccessCode":   c.Secret,
			"TravelTimeID": travelTimeID,
		}).
		Get(travelTimeURL)

	if err != nil {
		log.WithFields(log.Fields{
			"error":        err,
			"travelTimeID": travelTimeID,
		}).Error("Unable to lokup the travel time specified")
		return nil, err
	}

	return resp.Result().(*wsdotTravelTimeAPIResponse), nil
}

func (c *serviceClient) adapt(info *wsdotTravelTimeAPIResponse) (twomqtt.Event, error) {
	log.WithFields(log.Fields{
		"info": info,
	}).Debug("Adapting travel time information")

	obj := wsdotTravelTime{
		CurrentTime:  info.CurrentTime,
		Distance:     info.Distance,
		TravelTimeID: info.TravelTimeID,
	}
	event := twomqtt.Event{
		Type:    reflect.TypeOf(obj),
		Payload: obj,
	}

	log.Debug("Finished adapting time travel information")
	return event, nil
}
