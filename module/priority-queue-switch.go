package learningrobotics

import (
	"container/heap"
	"context"
	"errors"
	"slices"
	"strconv"
	"sync"
	"time"

	sw "go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	generic "go.viam.com/rdk/services/generic"
)

var (
	PriorityQueueSwitch = resource.NewModel("mattmacf", "learning-robotics", "priority-queue-switch")
	errUnimplemented    = errors.New("unimplemented")
)

func init() {
	resource.RegisterService(generic.API, PriorityQueueSwitch,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newPriorityQueueSwitchPriorityQueueSwitch,
		},
	)
}

type Config struct {
	SwitchName string `json:"switch_name"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
	if cfg.SwitchName == "" {
		return nil, nil, errors.New("switch_name is required")
	}
	return nil, nil, nil
}

type priorityQueueSwitchPriorityQueueSwitch struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()
	sw         sw.Switch
	pq         PriorityQueue
	mu         sync.Mutex
}

func newPriorityQueueSwitchPriorityQueueSwitch(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewPriorityQueueSwitch(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewPriorityQueueSwitch(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (resource.Resource, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	sw, err := sw.FromProvider(deps, conf.SwitchName)
	if err != nil {
		cancelFunc()
		return nil, err
	}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	s := &priorityQueueSwitchPriorityQueueSwitch{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		sw:         sw,
		pq:         pq,
		mu:         sync.Mutex{},
	}
	go s.drainPriorityQueue()
	return s, nil
}

func (s *priorityQueueSwitchPriorityQueueSwitch) Name() resource.Name {
	return s.name
}

func (s *priorityQueueSwitchPriorityQueueSwitch) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := cmd["get_length"]; ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		return map[string]any{"length": s.pq.Len()}, nil
	}

	label, ok := cmd["label"].(string)
	if !ok {
		return nil, errors.New("label is required")
	}
	priorityStr, ok := cmd["priority"].(string)
	if !ok {
		return nil, errors.New("priority is required")
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		return nil, errors.New("priority must be an integer")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, validLabels, err := s.sw.GetNumberOfPositions(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if !slices.Contains(validLabels, label) {
		return nil, errors.New("label is not valid")
	}
	position := slices.Index(validLabels, label)

	heap.Push(&s.pq, &CommandItem{position: position, priority: priority})

	return nil, nil
}

func (s *priorityQueueSwitchPriorityQueueSwitch) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}

func (s *priorityQueueSwitchPriorityQueueSwitch) drainPriorityQueue() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-s.cancelCtx.Done():
			return
		case <-ticker.C:
			s.mu.Lock()
			defer s.mu.Unlock()

			if s.pq.Len() > 0 {
				item := heap.Pop(&s.pq).(*CommandItem)
				s.sw.SetPosition(s.cancelCtx, uint32(item.position), map[string]interface{}{})
			}
			s.mu.Unlock()
		}
	}
}

type CommandItem struct {
	position int
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
