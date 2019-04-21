package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	resty "gopkg.in/resty.v1"
)

const (
	travelTimeURL = "http://www.wsdot.wa.gov/Traffic/api/TravelTimes/TravelTimesREST.svc/GetTravelTimeAsJson"
)

type client struct {
	observers map[observer]struct{}

	lookupInterval time.Duration
	travelTimes    map[string]string
	secret         string
}

func newClient(config *config) *client {
	c := client{
		observers: map[observer]struct{}{},

		lookupInterval: config.LookupInterval,
		secret:         config.Secret,
	}
	c.travelTimes = make(map[string]string, 0)

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

	return &c
}

func (c *client) run() {
	go c.loop(false)
}

func (c *client) register(l observer) {
	c.observers[l] = struct{}{}
}

func (c *client) publish(e event) {
	for o := range c.observers {
		o.receive(e)
	}
}

func (c *client) loop(once bool) {
	for {
		log.Print("Beginning lookup")
		for travelTimeID, travelTimeSlug := range c.travelTimes {
			if info, err := c.lookup(travelTimeID); err == nil {
				c.publish(event{
					version: 1,
					key:     travelTimeSlug,
					data:    *info,
				})
			} else {
				log.Print(err)
			}
		}
		log.Print("Ending lookup")

		if once {
			break
		}

		time.Sleep(c.lookupInterval)
	}
}

func (c *client) lookup(travelTimeID string) (*travelTime, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&travelTime{}).
		SetQueryParams(map[string]string{
			"AccessCode":   c.secret,
			"TravelTimeID": travelTimeID,
		}).
		Get(travelTimeURL)

	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("Unble to lookup the travel time for %s", travelTimeID)
	}

	return resp.Result().(*travelTime), nil
}
