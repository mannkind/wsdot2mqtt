package main

type wsdotTravelTime struct {
	AverageTime int     `json:"AverageTime"`
	CurrentTime int     `json:"CurrentTime"`
	Description string  `json:"Description"`
	Distance    float64 `json:"Distance"`
	EndPoint    struct {
		Description string  `json:"Description"`
		Direction   string  `json:"Direction"`
		Latitude    float64 `json:"Latitude"`
		Longitude   float64 `json:"Longitude"`
		MilePost    float64 `json:"MilePost"`
		RoadName    string  `json:"RoadName"`
	} `json:"EndPoint"`
	Name       string `json:"Name"`
	StartPoint struct {
		Description string  `json:"Description"`
		Direction   string  `json:"Direction"`
		Latitude    float64 `json:"Latitude"`
		Longitude   float64 `json:"Longitude"`
		MilePost    float64 `json:"MilePost"`
		RoadName    string  `json:"RoadName"`
	} `json:"StartPoint"`
	TimeUpdated  string `json:"TimeUpdated"`
	TravelTimeID int    `json:"TravelTimeID"`
}

type eventData struct {
	CurrentTime  int
	Distance     float64
	TravelTimeID int
}

type event struct {
	version int64
	key     string
	data    eventData
}

type observer interface {
	receiveState(event)
	receiveCommand(int64, event)
}

type publisher interface {
	register(observer)
}
