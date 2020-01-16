package main

import (
	"sensor"
	"sensor/mq"
)

func main() {
	mq.SensorMapping()
	sensor.SensorServiceStart()
}
