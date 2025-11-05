package learningrobotics

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	generic "go.viam.com/rdk/services/generic"
)

var (
	LightSwitch      = resource.NewModel("mattmacf", "learning-robotics", "light-switch")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterService(generic.API, LightSwitch,
		resource.Registration[resource.Resource, *LightSwitchConfig]{
			Constructor: newLearningRoboticsLightSwitch,
		},
	)
}

type LightSwitchConfig struct {
	LightOutputPin    string `json:"light_output_pin"`
	OnButtonInputPin  string `json:"on_button_input_pin"`
	OffButtonInputPin string `json:"off_button_input_pin"`
	BoardName         string `json:"board_name"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *LightSwitchConfig) Validate(path string) ([]string, []string, error) {
	if cfg.LightOutputPin == "" {
		return nil, nil, errors.New("light_output_pin is required")
	}
	if cfg.OnButtonInputPin == "" {
		return nil, nil, errors.New("on_button_input_pin is required")
	}
	if cfg.OffButtonInputPin == "" {
		return nil, nil, errors.New("off_button_input_pin is required")
	}
	if cfg.BoardName == "" {
		return nil, nil, errors.New("board_name is required")
	}
	return nil, nil, nil
}

type learningRoboticsLightSwitch struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *LightSwitchConfig

	cancelCtx  context.Context
	cancelFunc func()

	lightOutputPin    board.GPIOPin
	onButtonInputPin  board.GPIOPin
	offButtonInputPin board.GPIOPin
}

func newLearningRoboticsLightSwitch(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*LightSwitchConfig](rawConf)
	if err != nil {
		return nil, err
	}

	return NewLightSwitch(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewLightSwitch(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *LightSwitchConfig, logger logging.Logger) (resource.Resource, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	board, err := board.FromProvider(deps, conf.BoardName)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	lightOutputPin, err := board.GPIOPinByName(conf.LightOutputPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	onButtonInputPin, err := board.GPIOPinByName(conf.OnButtonInputPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	offButtonInputPin, err := board.GPIOPinByName(conf.OffButtonInputPin)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	s := &learningRoboticsLightSwitch{
		name:              name,
		logger:            logger,
		cfg:               conf,
		cancelCtx:         cancelCtx,
		cancelFunc:        cancelFunc,
		lightOutputPin:    lightOutputPin,
		onButtonInputPin:  onButtonInputPin,
		offButtonInputPin: offButtonInputPin,
	}
	go s.run(cancelCtx)
	return s, nil
}

func (s *learningRoboticsLightSwitch) Name() resource.Name {
	return s.name
}

func (s *learningRoboticsLightSwitch) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *learningRoboticsLightSwitch) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}

func (s *learningRoboticsLightSwitch) run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second / 30) // 30 frames per second
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			lightOnValue, err := s.onButtonInputPin.Get(ctx, map[string]interface{}{})
			if err != nil {
				return err
			}
			if !lightOnValue {
				s.lightOutputPin.Set(ctx, true, map[string]interface{}{})
			}

			lightOffValue, err := s.offButtonInputPin.Get(ctx, map[string]interface{}{})
			if err != nil {
				return err
			}
			if !lightOffValue {
				s.lightOutputPin.Set(ctx, false, map[string]interface{}{})
			}
		}
	}
}
