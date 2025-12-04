package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/robot/client"
)

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

func TestEventSystem() {
	fmt.Println("Testing event system...")
	mq := NewMessageQueue(10)
	mq.Subscribe(func(message EventMessage) {
		fmt.Println("Received message from subscriber 1:", message)
	})
	mq.Subscribe(func(message EventMessage) {
		fmt.Println("Received message from subscriber 2:", message)
	})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "hello"})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "world"})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "foo"})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "bar"})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "baz"})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "qux"})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "quux"})
	fmt.Println("Adding subscriber 3")
	mq.Subscribe(func(message EventMessage) {
		fmt.Println("Received message from subscriber 3:", message)
	})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "corge"})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "grault"})
	time.Sleep(time.Second * 1)
	mq.Publish(EventMessage{topic: "test", data: "the quick brown fox jumps over the lazy dog"})

	// Block until mq.ch is empty
	for {
		mq.mu.Lock()
		lenCh := len(mq.ch)
		mq.mu.Unlock()
		if lenCh == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

const RED_LED_PIN_EVENT_SYSTEM = "8"
const GREEN_LED_PIN_EVENT_SYSTEM = "10"
const BLUE_LED_PIN_EVENT_SYSTEM = "12"
const BUZZER_PIN_EVENT_SYSTEM = "32"
const TRIGGER_PIN_EVENT_SYSTEM = "35"
const ECHO_PIN_EVENT_SYSTEM = "37"

func EventSystem(machine *client.RobotClient) error {
	ticksChan := make(chan board.Tick, 2)
	defer close(ticksChan)

	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pi, err := board.FromProvider(machine, "pi")
	if err != nil {
		return err
	}

	red_pin, err := pi.GPIOPinByName(RED_LED_PIN_EVENT_SYSTEM)
	if err != nil {
		return err
	}

	blue_pin, err := pi.GPIOPinByName(BLUE_LED_PIN_EVENT_SYSTEM)
	if err != nil {
		return err
	}
	green_pin, err := pi.GPIOPinByName(GREEN_LED_PIN_EVENT_SYSTEM)
	if err != nil {
		return err
	}

	buzzer_pin, err := pi.GPIOPinByName(BUZZER_PIN_EVENT_SYSTEM)
	if err != nil {
		return err
	}

	triggerPin, err := pi.GPIOPinByName(TRIGGER_PIN_EVENT_SYSTEM)
	if err != nil {
		return err
	}
	echoPin, err := pi.DigitalInterruptByName(ECHO_PIN_EVENT_SYSTEM)
	if err != nil {
		return err
	}

	mq := NewMessageQueue(10)

	mq.Subscribe(func(message EventMessage) {
		distance := message.data.(float64)
		if distance < 0.3 {
			makeRed(red_pin, blue_pin, green_pin)
		} else {
			makeGreen(red_pin, blue_pin, green_pin)
		}
	})

	mq.Subscribe(func(message EventMessage) {
		distance := message.data.(float64)

		if distance < 0.1 {
			buzzer_pin.SetPWM(context.Background(), 0.05, map[string]interface{}{})
			buzzer_pin.SetPWMFreq(context.Background(), 1000, map[string]interface{}{})
		} else if distance <= 0.4 {
			buzzer_pin.SetPWM(context.Background(), 0.1, map[string]interface{}{})
			buzzer_pin.SetPWMFreq(context.Background(), 800, map[string]interface{}{})
		} else if distance <= 0.7 {
			buzzer_pin.SetPWM(context.Background(), 0.2, map[string]interface{}{})
			buzzer_pin.SetPWMFreq(context.Background(), 500, map[string]interface{}{})
		} else {
			buzzer_pin.SetPWM(context.Background(), 0, map[string]interface{}{})
		}
	})

	pi.StreamTicks(cancelCtx, []board.DigitalInterrupt{echoPin}, ticksChan, map[string]interface{}{})

	ticker := time.NewTicker(time.Millisecond * 100) // ( recommended of 60 ms between readings)
	defer ticker.Stop()

	for range ticker.C {
		distMeters, err := getReading(cancelCtx, triggerPin, ticksChan)
		if err != nil {
			return err
		}
		mq.Publish(EventMessage{topic: "distance", data: distMeters})
		fmt.Printf("Distance to nearest object: %f meters\n", distMeters)
	}
	return nil
}
