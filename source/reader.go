package source

import (
	"fmt"

	"github.com/mannkind/wsdot2mqtt/shared"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// Reader is for reading a shared representation out of a source system
type Reader struct {
	opts     Opts
	outgoing chan<- shared.Representation
	service  *Service
}

// NewReader creates a new Reader for reading a shared representation out of a source system
func NewReader(opts Opts, outgoing chan<- shared.Representation, service *Service) *Reader {
	c := Reader{
		opts:     opts,
		outgoing: outgoing,
		service:  service,
	}

	service.SetSecret(opts.Secret)

	return &c
}

// Run starts the Reader
func (c *Reader) Run() {
	// Log service settings
	c.logSettings()

	// Run immediately
	c.poll()

	// Schedule additional runs
	sched := cron.New()
	sched.AddFunc(fmt.Sprintf("@every %s", c.opts.LookupInterval), c.poll)
	sched.Start()
}

func (c *Reader) logSettings() {
	// Log the current settings
	log.WithFields(log.Fields{
		"WSDOT.LookupInterval": c.opts.LookupInterval,
		"WSDOT.TravelMapping":  c.opts.TravelTimeMapping,
	}).Info("Service Environmental Settings")
}

// poll the source system, adapt source system responses to the share representation, output data onto a channnel
func (c *Reader) poll() {
	log.Info("Polling")
	for travelTimeID := range c.opts.TravelTimeMapping {
		info, err := c.service.lookup(travelTimeID)
		if err != nil {
			continue
		}

		c.outgoing <- c.adapt(info)
	}

	log.WithFields(log.Fields{
		"sleep": c.opts.LookupInterval,
	}).Info("Finished polling; sleeping")
}

// adapt incoming value(s) to the shared representation
func (c *Reader) adapt(info *serviceRepresentation) shared.Representation {
	return shared.Representation{
		CurrentTime:  info.CurrentTime,
		Distance:     info.Distance,
		TravelTimeID: info.TravelTimeID,
	}
}
