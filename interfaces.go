package main

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
	receive(event)
}

type publisher interface {
	register(observer)
	publish(event)
}
