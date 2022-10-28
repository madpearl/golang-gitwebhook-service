package validator

import (
	"fmt"
	"os"
	"testing"

	"github.com/microlib/simple"
)

func TestEnvars(t *testing.T) {
	logger := &simple.Logger{Level: "info"}

	t.Run("ValidateEnvars : should fail", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "")
		err := ValidateEnvars(logger)
		if err == nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with no error - got (%v) wanted (%v)", "ValidateEnvars", err, nil))
		}
	})

	t.Run("ValidateEnvars : should pass", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "info")
		os.Setenv("PR_OPENED_URL", "localhost")
		os.Setenv("PR_MERGED_URL", "localhost")
		os.Setenv("PRERELEASED_URL", "localhost")
		os.Setenv("RELEASED_URL", "localhost")
		os.Setenv("VERSION", "1.0.3")
		os.Setenv("WEBHOOK_SECRET", "ewqewqe")
		os.Setenv("REPO_MAPPING", "test")
		os.Setenv("NAME", "test")
		err := ValidateEnvars(logger)
		if err != nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with error - got (%v) wanted (%v)", "ValidateEnvars", err, nil))
		}
	})

}
