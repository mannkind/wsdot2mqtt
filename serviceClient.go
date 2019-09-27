package main

import (
	"fmt"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

const (
	travelTimeURL = "https://www.wsdot.wa.gov/Traffic/api/TravelTimes/TravelTimesREST.svc/GetTravelTimeAsJson"
)

type serviceClient struct {
	serviceClientConfig
	stateUpdateChan stateChannel
}

func newServiceClient(serviceClientCfg serviceClientConfig, stateUpdateChan stateChannel) *serviceClient {
	c := serviceClient{
		serviceClientConfig: serviceClientCfg,
		stateUpdateChan:     stateUpdateChan,
	}

	log.WithFields(log.Fields{
		"WSDOT.LookupInterval": c.LookupInterval,
		"WSDOT.TravelMapping":  c.TravelTimeMapping,
	}).Info("Service Environmental Settings")

	return &c
}

func (c *serviceClient) run() {
	// Run immediately
	go c.loop()

	// Schedule additional runs
	sched := cron.New()
	sched.AddFunc(fmt.Sprintf("@every %s", c.LookupInterval), c.loop)
	sched.Start()
}

func (c *serviceClient) loop() {
	log.Info("Looping")
	for travelTimeID := range c.TravelTimeMapping {
		info, err := c.lookup(travelTimeID)
		if err != nil {
			continue
		}

		obj, err := c.adapt(info)
		if err != nil {
			continue
		}

		c.stateUpdateChan <- obj
	}

	log.WithFields(log.Fields{
		"sleep": c.LookupInterval,
	}).Info("Finished looping; sleeping")
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

func (c *serviceClient) adapt(info *wsdotTravelTimeAPIResponse) (wsdotTravelTime, error) {
	log.WithFields(log.Fields{
		"info": info,
	}).Debug("Adapting travel time information")

	obj := wsdotTravelTime{
		CurrentTime:  info.CurrentTime,
		Distance:     info.Distance,
		TravelTimeID: info.TravelTimeID,
	}

	log.Debug("Finished adapting time travel information")
	return obj, nil
}
