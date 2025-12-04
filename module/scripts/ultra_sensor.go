package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/robot/client"
)

const TRIGGER_PIN = "8"
const ECHO_PIN = "10"

func UltraSensorRead(machine *client.RobotClient) error {
	ticksChan := make(chan board.Tick, 2)
	defer close(ticksChan)

	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pi, err := board.FromProvider(machine, "pi")
	if err != nil {
		return err
	}
	triggerPin, err := pi.GPIOPinByName(TRIGGER_PIN)
	if err != nil {
		return err
	}
	echoPin, err := pi.DigitalInterruptByName(ECHO_PIN)
	if err != nil {
		return err
	}

	pi.StreamTicks(cancelCtx, []board.DigitalInterrupt{echoPin}, ticksChan, map[string]interface{}{})

	ticker := time.NewTicker(time.Second) // ( recommended of 60 ms between readings)
	defer ticker.Stop()

	for range ticker.C {
		distMeters, err := getReading(cancelCtx, triggerPin, ticksChan)
		if err != nil {
			return err
		}
		fmt.Printf("Distance to nearest object: %f meters\n", distMeters)
	}
	return nil
}

func getReading(ctx context.Context, triggerPin board.GPIOPin, ticksChan <-chan board.Tick) (float64, error) {

	// set trigger pin low
	triggerPin.Set(ctx, false, map[string]interface{}{})

	// we send a high and a low to the trigger pin 10 microseconds
	// apart to signal the sensor to begin sending the sonic pulse
	if err := triggerPin.Set(ctx, true, map[string]interface{}{}); err != nil {
		return 0, err
	}
	time.Sleep(time.Microsecond * 10)
	if err := triggerPin.Set(ctx, false, map[string]interface{}{}); err != nil {
		return 0, err
	}

	// the first signal from the interrupt indicates that the sonic
	// pulse has been sent and the second indicates that the echo has been received
	var tick board.Tick
	ticks := make([]board.Tick, 2)

	for i := range 2 {
		var signalStr string
		if i == 0 {
			signalStr = "sound pulse was emitted"
		} else {
			signalStr = "echo was received"
		}
		select {
		case tick = <-ticksChan:
			ticks[i] = tick
		case <-ctx.Done():
			fmt.Printf("Context cancelled while waiting for signal that %s\n", signalStr)
			return 0, ctx.Err()
		}
	}

	timeEmitted := ticks[0].TimestampNanosec
	timeReceived := ticks[1].TimestampNanosec
	// we calculate the distance to the nearest object based
	// on the time interval between the sound and its echo
	// and the speed of sound (343 m/s)
	secondsElapsed := float64(timeReceived-timeEmitted) / math.Pow10(9)
	distMeters := secondsElapsed * 343 / 2
	return distMeters, nil
}
