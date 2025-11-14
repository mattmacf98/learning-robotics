package main

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/robot/client"
)

type Color string

const (
	Red   Color = "Red"
	Green Color = "Green"
	Blue  Color = "Blue"
)

const RED_LED_PIN_PRIORITY = "8"
const GREEN_LED_PIN_PRIORITY = "10"
const BLUE_LED_PIN_PRIORITY = "12"

type CommandItem struct {
	color    Color
	priority int
	index    int
}

type PriorityQueue []*CommandItem

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*CommandItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func RgbPriorityQueue(machine *client.RobotClient) error {
	light_commands_pq := make(PriorityQueue, 0)
	heap.Init(&light_commands_pq)

	pi, err := board.FromProvider(machine, "pi")
	if err != nil {
		return err
	}

	red_pin, err := pi.GPIOPinByName(RED_LED_PIN_PRIORITY)
	if err != nil {
		return err
	}
	blue_pin, err := pi.GPIOPinByName(BLUE_LED_PIN_PRIORITY)
	if err != nil {
		return err
	}
	green_pin, err := pi.GPIOPinByName(GREEN_LED_PIN_PRIORITY)
	if err != nil {
		return err
	}

	items := []*CommandItem{
		{color: "red", priority: 2},
		{color: "blue", priority: 5},
		{color: "green", priority: 1},
		{color: "red", priority: 3},
	}

	for _, item := range items {
		heap.Push(&light_commands_pq, item)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for light_commands_pq.Len() > 0 {
			item := heap.Pop(&light_commands_pq).(*CommandItem)
			switch item.color {
			case "red":
				fmt.Println("making red")
				makeRedPriority(red_pin, blue_pin, green_pin)
			case "green":
				fmt.Println("making green")
				makeGreenPriority(red_pin, blue_pin, green_pin)
			case "blue":
				fmt.Println("making blue")
				makeBluePriority(red_pin, blue_pin, green_pin)
			}
			time.Sleep(time.Millisecond * 1000)
			turnOffPriority(red_pin, blue_pin, green_pin)
			time.Sleep(time.Millisecond * 1000)
		}
	}()

	wg.Wait()
	fmt.Println("done")
	return nil
}

func makeRedPriority(r board.GPIOPin, b board.GPIOPin, g board.GPIOPin) {
	r.Set(context.Background(), true, map[string]interface{}{})
	b.Set(context.Background(), false, map[string]interface{}{})
	g.Set(context.Background(), false, map[string]interface{}{})
}

func makeGreenPriority(r board.GPIOPin, b board.GPIOPin, g board.GPIOPin) {
	r.Set(context.Background(), false, map[string]interface{}{})
	b.Set(context.Background(), false, map[string]interface{}{})
	g.Set(context.Background(), true, map[string]interface{}{})
}

func makeBluePriority(r board.GPIOPin, b board.GPIOPin, g board.GPIOPin) {
	r.Set(context.Background(), false, map[string]interface{}{})
	b.Set(context.Background(), true, map[string]interface{}{})
	g.Set(context.Background(), false, map[string]interface{}{})
}

func turnOffPriority(r board.GPIOPin, b board.GPIOPin, g board.GPIOPin) {
	r.Set(context.Background(), false, map[string]interface{}{})
	b.Set(context.Background(), false, map[string]interface{}{})
	g.Set(context.Background(), false, map[string]interface{}{})
}
