package main

type stateChannel = chan wsdotTravelTime

func newStateChannel() stateChannel {
	return make(stateChannel, 100)
}
