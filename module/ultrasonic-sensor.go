package learningrobotics

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"go.viam.com/rdk/components/board"
	sensor "go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var (
	UltrasonicSensor = resource.NewModel("mattmacf", "learning-robotics", "ultrasonic-sensor")
)

func init() {
	resource.RegisterComponent(sensor.API, UltrasonicSensor,
		resource.Registration[sensor.Sensor, *UltrasonicSensorConfig]{
			Constructor: newUltraSensorUltrasonicSensor,
		},
	)
}

type UltrasonicSensorConfig struct {
	TriggerPin    string `json:"trigger_pin"`
	EchoInterrupt string `json:"echo_interrupt_pin"`
	BoardName     string `json:"board_name"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *UltrasonicSensorConfig) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
	if cfg.TriggerPin == "" {
		return nil, nil, errors.New("trigger_pin is required")
	}
	if cfg.EchoInterrupt == "" {
		return nil, nil, errors.New("echo_interrupt_pin is required")
	}
	if cfg.BoardName == "" {
		return nil, nil, errors.New("board_name is required")
	}
	return []string{cfg.BoardName}, nil, nil
}

type ultraSensorUltrasonicSensor struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *UltrasonicSensorConfig

	cancelCtx  context.Context
	cancelFunc func()

	triggerPin    board.GPIOPin
	echoInterrupt board.DigitalInterrupt
	ticksChan     chan board.Tick
}

func newUltraSensorUltrasonicSensor(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	conf, err := resource.NativeConfig[*UltrasonicSensorConfig](rawConf)
	if err != nil {
		return nil, err
	}

	return NewUltrasonicSensor(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewUltrasonicSensor(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *UltrasonicSensorConfig, logger logging.Logger) (sensor.Sensor, error) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	ticksChan := make(chan board.Tick, 2)

	piBoard, err := board.FromProvider(deps, conf.BoardName)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	triggerPin, err := piBoard.GPIOPinByName(conf.TriggerPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	echoInterrupt, err := piBoard.DigitalInterruptByName(conf.EchoInterrupt)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	s := &ultraSensorUltrasonicSensor{
		name:          name,
		logger:        logger,
		cfg:           conf,
		cancelCtx:     cancelCtx,
		cancelFunc:    cancelFunc,
		triggerPin:    triggerPin,
		echoInterrupt: echoInterrupt,
		ticksChan:     ticksChan,
	}

	piBoard.StreamTicks(cancelCtx, []board.DigitalInterrupt{echoInterrupt}, ticksChan, map[string]interface{}{})
	return s, nil
}

func (s *ultraSensorUltrasonicSensor) Name() resource.Name {
	return s.name
}

func (s *ultraSensorUltrasonicSensor) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	// set trigger pin low
	if err := s.triggerPin.Set(ctx, false, map[string]interface{}{}); err != nil {
		return nil, err
	}

	// we send a high and a low to the trigger pin 10 microseconds
	// apart to signal the sensor to begin sending the sonic pulse
	if err := s.triggerPin.Set(ctx, true, map[string]interface{}{}); err != nil {
		return nil, err
	}
	time.Sleep(time.Microsecond * 10)
	if err := s.triggerPin.Set(ctx, false, map[string]interface{}{}); err != nil {
		return nil, err
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
		case tick = <-s.ticksChan:
			ticks[i] = tick
		case <-ctx.Done():
			fmt.Printf("Context cancelled while waiting for signal that %s\n", signalStr)
			return nil, ctx.Err()
		}
	}

	timeEmitted := ticks[0].TimestampNanosec
	timeReceived := ticks[1].TimestampNanosec
	// we calculate the distance to the nearest object based
	// on the time interval between the sound and its echo
	// and the speed of sound (343 m/s)
	secondsElapsed := float64(timeReceived-timeEmitted) / math.Pow10(9)
	distMeters := secondsElapsed * 343 / 2
	return map[string]interface{}{"distance": distMeters}, nil
}

func (s *ultraSensorUltrasonicSensor) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *ultraSensorUltrasonicSensor) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
