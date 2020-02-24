package source

// serviceRepresentation is a representation of the data returned from the source system
type serviceRepresentation struct {
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
