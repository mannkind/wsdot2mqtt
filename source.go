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

type source struct {
	config   sourceOpts
	outgoing chan<- sourceRep
}

func newSource(config sourceOpts, outgoing chan<- sourceRep) *source {
	c := source{
		config:   config,
		outgoing: outgoing,
	}

	return &c
}

func (c *source) run() {
	// Log service settings
	c.logSettings()

	// Run immediately
	c.poll()

	// Schedule additional runs
	sched := cron.New()
	sched.AddFunc(fmt.Sprintf("@every %s", c.config.LookupInterval), c.poll)
	sched.Start()
}

func (c *source) logSettings() {
	// Log the current settings
	log.WithFields(log.Fields{
		"WSDOT.LookupInterval": c.config.LookupInterval,
		"WSDOT.TravelMapping":  c.config.TravelTimeMapping,
	}).Info("Service Environmental Settings")
}

func (c *source) poll() {
	log.Info("Polling")
	for travelTimeID := range c.config.TravelTimeMapping {
		info, err := c.lookup(travelTimeID)
		if err != nil {
			continue
		}

		c.outgoing <- c.adapt(info)
	}

	log.WithFields(log.Fields{
		"sleep": c.config.LookupInterval,
	}).Info("Finished polling; sleeping")
}

func (c *source) lookup(travelTimeID string) (*sourceResponse, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&sourceResponse{}).
		SetQueryParams(map[string]string{
			"AccessCode":   c.config.Secret,
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

	return resp.Result().(*sourceResponse), nil
}

func (c *source) adapt(info *sourceResponse) sourceRep {
	return sourceRep{
		CurrentTime:  info.CurrentTime,
		Distance:     info.Distance,
		TravelTimeID: info.TravelTimeID,
	}
}
