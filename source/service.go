package source

import (
	log "github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

// Service is for reading a directly from a source system
type Service struct {
	Secret string
}

// NewService creates a new Service for reading a directly from a source system
func NewService() *Service {
	c := Service{}

	return &c
}

// SetSecret sets the required options to access the source system
func (c *Service) SetSecret(secret string) {
	c.Secret = secret
}

// lookup data from the source system
func (c *Service) lookup(travelTimeID string) (*serviceRepresentation, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&serviceRepresentation{}).
		SetQueryParams(map[string]string{
			"AccessCode":   c.Secret,
			"TravelTimeID": travelTimeID,
		}).
		Get("https://www.wsdot.wa.gov/Traffic/api/TravelTimes/TravelTimesREST.svc/GetTravelTimeAsJson")

	if err != nil {
		log.WithFields(log.Fields{
			"error":        err,
			"travelTimeID": travelTimeID,
		}).Error("Unable to lokup the travel time specified")
		return nil, err
	}

	return resp.Result().(*serviceRepresentation), nil
}
