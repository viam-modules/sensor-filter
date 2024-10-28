package sensor_filter

// Test Validate .. both an empty config and a valid one, maybe also bad but not empty
// Test DoCommand.. make sure it returns a map with the right type

import (
	"testing"

	"go.viam.com/test"
)

func TestValidate(t *testing.T) {
	// empty
	cfg := &Config{"", "", nil}
	_, err := cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)
	
}