package sensor_filter

// Test Validate .. both an empty config and a valid one, maybe also bad but not empty
// Test DoCommand.. make sure it returns a map with the right type

import (
	"context"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/testutils/inject"
	"go.viam.com/utils/apputils"
	"testing"

	"go.viam.com/test"
)

func TestValidate(t *testing.T) {
	// empty
	cfg := &Config{"", "", nil}
	_, err := cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)

	// no SensorName
	cfg = &Config{"", "temperature", nil}
	_, err = cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, "error validating")

	// no reading
	cfg = &Config{"mySensor", "", nil}
	_, err = cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, err.Error(), test.ShouldContainSubstring, "error validating")

	// good vibes
	cond := []*apputils.Eval{{Operator: "lt", Value: 2}}
	cfg = &Config{"mySensor", "temperature", cond}
	_, err = cfg.Validate("")
	test.That(t, err, test.ShouldBeNil)

}

func TestDo(t *testing.T) {
	sens := inject.NewSensor("test")
	sens.ReadingsFunc = func(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
		out := make(map[string]interface{})
		out["distance"] = 2.718281828459
		return out, nil
	}

	cond := []*apputils.Eval{{Operator: "lt", Value: 2}}
	cfg := &Config{"test", "distance", cond}
	rcfg := resource.Config{Name: "mySF", ConvertedAttributes: cfg}
	_, err := cfg.Validate("")
	test.That(t, err, test.ShouldBeNil)

	ctx := context.Background()
	logger := logging.NewTestLogger(t)
	resourceMap := map[resource.Name]resource.Resource{
		sensor.Named("test"): sens,
	}

	sf, err := newSensorFilter(ctx, resourceMap, rcfg, logger)
	test.That(t, err, test.ShouldBeNil)

	err = sf.Reconfigure(ctx, resourceMap, rcfg)
	test.That(t, err, test.ShouldBeNil)

	out, err := sf.DoCommand(ctx, nil)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, out["result"], test.ShouldNotBeNil)
	test.That(t, out["result"], test.ShouldEqual, false) // 2.718 > 2
}
