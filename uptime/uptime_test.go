package uptime

import (
	"reflect"
	"testing"

	"github.com/bronger/jfy/lib"
)

func TestHandle(t *testing.T) {
	cases := []struct {
		input    []byte
		expected map[string]any
	}{
		{[]byte(" 20:00:50 up  1:01,  4 users,  load average: 0.72, 0.95, 0.95\n"),
			map[string]any{"days": 0, "hour": 20, "hours": 1, "load1": 0.72, "load15": 0.95, "load5": 0.95,
				"minute": 0, "minutes": 1, "second": 50, "users": 4}},
	}
	for _, tc := range cases {
		output, errors, err :=
			Handle(lib.SettingsType{}, tc.input, []byte{}, "uptime")
		if err != nil {
			t.Errorf("Got internal error '%v'", err)
		}
		if errors != nil {
			t.Errorf("Expected no error output and got %v", errors)
		}
		if !reflect.DeepEqual(output, tc.expected) {
			t.Errorf("Output was wrongly %v", output)
		}
	}
}
