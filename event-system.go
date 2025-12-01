package learningrobotics

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.viam.com/rdk/components/board"
	sensor "go.viam.com/rdk/components/sensor"
	sw "go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	generic "go.viam.com/rdk/services/generic"
)

var (
	EventSystem = resource.NewModel("mattmacf", "learning-robotics", "event-system")
)

func init() {
	resource.RegisterService(generic.API, EventSystem,
		resource.Registration[resource.Resource, *EventSystemConfig]{
			Constructor: newEventSystemEventSystem,
		},
	)
}

type EventSystemConfig struct {
	UltrasonicSensorName string `json:"ultrasonic_sensor_name"`
	RGBSwitchName        string `json:"rgb_switch_name"`
	BuzzerPin            string `json:"buzzer_pin"`
	BoardName            string `json:"board_name"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *EventSystemConfig) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
	if cfg.UltrasonicSensorName == "" {
		return nil, nil, errors.New("ultrasonic_sensor_name is required")
	}
	if cfg.RGBSwitchName == "" {
		return nil, nil, errors.New("rgb_switch_name is required")
	}
	if cfg.BuzzerPin == "" {
		return nil, nil, errors.New("buzzer_pin is required")
	}
	if cfg.BoardName == "" {
		return nil, nil, errors.New("board_name is required")
	}
	return nil, nil, nil
}

type eventSystemEventSystem struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *EventSystemConfig

	cancelCtx  context.Context
	cancelFunc func()

	buzzerPin        board.GPIOPin
	rgbSwitch        sw.Switch
	ultrasonicSensor sensor.Sensor
	mq               *MessageQueue
}

func newEventSystemEventSystem(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*EventSystemConfig](rawConf)
	if err != nil {
		return nil, err
	}

	return NewEventSystem(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewEventSystem(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *EventSystemConfig, logger logging.Logger) (resource.Resource, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	board, err := board.FromProvider(deps, conf.BoardName)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	buzzerPin, err := board.GPIOPinByName(conf.BuzzerPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	rgbSwitch, err := sw.FromProvider(deps, conf.RGBSwitchName)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	ultrasonicSensor, err := sensor.FromProvider(deps, conf.UltrasonicSensorName)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	mq := NewMessageQueue(10)
	s := &eventSystemEventSystem{
		name:             name,
		logger:           logger,
		cfg:              conf,
		cancelCtx:        cancelCtx,
		cancelFunc:       cancelFunc,
		buzzerPin:        buzzerPin,
		rgbSwitch:        rgbSwitch,
		ultrasonicSensor: ultrasonicSensor,
		mq:               mq,
	}

	mq.Subscribe(func(message EventMessage) {
		distance := message.data.(float64)
		if distance < 0.3 {
			// Red
			s.rgbSwitch.SetPosition(context.Background(), 1, map[string]interface{}{})
		} else {
			// Green
			s.rgbSwitch.SetPosition(context.Background(), 2, map[string]interface{}{})
		}
	})

	mq.Subscribe(func(message EventMessage) {
		distance := message.data.(float64)

		if distance < 0.1 {
			s.buzzerPin.SetPWM(context.Background(), 0.05, map[string]interface{}{})
			s.buzzerPin.SetPWMFreq(context.Background(), 1000, map[string]interface{}{})
		} else if distance <= 0.4 {
			s.buzzerPin.SetPWM(context.Background(), 0.1, map[string]interface{}{})
			s.buzzerPin.SetPWMFreq(context.Background(), 800, map[string]interface{}{})
		} else if distance <= 0.7 {
			s.buzzerPin.SetPWM(context.Background(), 0.2, map[string]interface{}{})
			s.buzzerPin.SetPWMFreq(context.Background(), 500, map[string]interface{}{})
		} else {
			s.buzzerPin.SetPWM(context.Background(), 0, map[string]interface{}{})
		}
	})

	go s.pollUltrasonicSensor(cancelCtx)

	return s, nil
}

func (s *eventSystemEventSystem) Name() resource.Name {
	return s.name
}

func (s *eventSystemEventSystem) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *eventSystemEventSystem) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}

func (s *eventSystemEventSystem) pollUltrasonicSensor(cancelCtx context.Context) {
	ticker := time.NewTicker(time.Millisecond * 100) // ( recommended of 60 ms between readings)
	defer ticker.Stop()

	for range ticker.C {
		readings, err := s.ultrasonicSensor.Readings(cancelCtx, map[string]interface{}{})
		if err != nil {
			continue
		}
		distance := readings["distance"].(float64)
		s.mq.Publish(EventMessage{topic: "distance", data: distance})
	}
}

type EventMessage struct {
	topic string
	data  any
}

type MessageQueue struct {
	ch          chan EventMessage
	subscribers []func(EventMessage)
	mu          sync.Mutex
}

func NewMessageQueue(bufferSize int) *MessageQueue {
	mq := &MessageQueue{
		ch:          make(chan EventMessage, bufferSize),
		subscribers: make([]func(EventMessage), 0),
		mu:          sync.Mutex{},
	}

	go mq.start()
	return mq
}

func (mq *MessageQueue) Subscribe(handler func(EventMessage)) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.subscribers = append(mq.subscribers, handler)
}

func (mq *MessageQueue) Publish(message EventMessage) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.ch <- message
}

func (mq *MessageQueue) start() {
	for message := range mq.ch {
		mq.mu.Lock()
		for _, subscriber := range mq.subscribers {
			go subscriber(message)
		}
		mq.mu.Unlock()
	}
}
