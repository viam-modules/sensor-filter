// Package sensor_filter implements a generic service.
package sensor_filter

import (
	"context"
	"github.com/pkg/errors"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/services/generic"
	"go.viam.com/utils"
	"go.viam.com/utils/apputils"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

// Here is where we define your new model's colon-delimited-triplet (viam-labs:go-module-templates-servo:customservo)
// viam-labs = namespace, go-module-templates-servo = repo-name, customservo = model name.
// TODO: Change model namespace, family (often the repo-name), and model. For more information see https://docs.viam.com/registry/create/#name-your-new-resource-model
var (
	Model          = resource.NewModel("viam", "generic", "sensor-filter")
	validOperators = []apputils.EvalOperator{
		apputils.Equal, apputils.NotEqual, apputils.LessThan, apputils.GreaterThan, apputils.LessThanOrEqual, apputils.GreaterThanOrEqual, apputils.Regex,
	}
)

func init() {
	resource.RegisterService(generic.API, Model,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newSensorFilter,
		},
	)
}

// Sensor filter service configuration. Each of the three fields are required.
type Config struct {
	SensorName string           `json:"sensor_name"`
	Reading    string           `json:"reading"`
	Conditions []*apputils.Eval `json:"conditions"`
}

// Validate validates the config and returns implicit dependencies.
func (cfg *Config) Validate(path string) ([]string, error) {

	if cfg.SensorName == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "sensor_name")
	}

	if cfg.Reading == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "reading")
	}

	for _, eval := range cfg.Conditions {
		if !eval.Operator.IsValidOperator(validOperators...) {
			return nil, errors.Errorf("cannot use %s operator", eval.Operator.ToReadableString())
		}
	}

	return []string{cfg.SensorName}, nil
}

// Constructor that creates and returns a new sensor filter service.
func newSensorFilter(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	// This takes the generic resource.Config passed down from the parent and converts it to the
	// model-specific (aka "native") Config structure defined above, making it easier to directly access attributes.
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	// Create a cancelable context
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &sensorFilter{
		name:       rawConf.ResourceName(),
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}

	// The Reconfigure() method changes the values on the service based on the attributes in the component config
	if err := s.Reconfigure(ctx, deps, rawConf); err != nil {
		logger.Error("Error configuring module with ", rawConf)
		return nil, err
	}

	return s, nil
}

// Methods in this service will be on this struct
type sensorFilter struct {
	name   resource.Name
	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	sensorName string
	sensor     sensor.Sensor
	reading    string
	condition  []*apputils.Eval
}

func (s *sensorFilter) Name() resource.Name {
	return s.name
}

// Reconfigures the model. Most models can be reconfigured in place without needing to rebuild. If you need to instead create a new instance of the servo, throw a NewMustBuildError.
func (s *sensorFilter) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	sfConfig, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		s.logger.Warn("Error reconfiguring module with ", err)
		return err
	}
	s.name = conf.ResourceName()

	// Make sure we can get the sensor
	s.sensorName = sfConfig.SensorName
	s.sensor, err = sensor.FromDependencies(deps, sfConfig.SensorName)
	if err != nil {
		return errors.Wrapf(err, "could not find sensor called %v ", s.sensorName)
	}

	// Get everything else
	s.condition = sfConfig.Conditions
	s.reading = sfConfig.Reading

	return nil
}

// The DoCommand for this module will evaluate a sensor reading against EVERY condition set up in the
// service config. The "result" will be true if and only if all conditions are met. Otherwise, it's false.
func (s *sensorFilter) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	// Read from sensor
	sensorOut, err := s.sensor.Readings(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Evaluate sensor reading against each condition
	// We want the logical AND of them all
	var finalAns bool
	for i, eval := range s.condition {
		ans, err := eval.Operator.Evaluate(sensorOut[s.reading], eval.Value)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		if i == 0 {
			finalAns = ans
		}
		finalAns = finalAns && ans
	}

	// Stuff result in a map and serve
	out := make(map[string]interface{})
	out["result"] = finalAns
	return out, err

}

// Close closes the underlying generic.
func (s *sensorFilter) Close(ctx context.Context) error {
	// NOT closing the sensor just b/c this service closes.
	s.cancelFunc()
	return nil
}
