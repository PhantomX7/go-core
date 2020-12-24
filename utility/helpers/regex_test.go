package helpers_test

import (
	"fmt"
	"testing"

	"github.com/PhantomX7/go-core/utility/helpers"
)

func TestCheckIsEmail(t *testing.T) {
	var checkIsEmailTests = []struct {
		input          string
		expectedOutput bool
	}{
		{"test@gmail.com", true},
		{"test", false},
		{"test@", false},
		{"test@co.id", true},
	}

	for _, tt := range checkIsEmailTests {
		t.Run(fmt.Sprint("check is email test input: ", tt.input), func(t *testing.T) {
			got := helpers.CheckIsEmail(tt.input)
			if got != tt.expectedOutput {
				t.Errorf("CheckIsEmails(%s) got %v, want %v", tt.input, got, tt.expectedOutput)
			}
		})
	}
}
