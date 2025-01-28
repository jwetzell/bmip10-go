package bmip10

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBMIP10Setup(t *testing.T) {
	var tests = []struct {
		description string
		sampleBits  uint32
		codedBits   uint32
		expected    Config
	}{
		{
			description: "10bit samples 8bit code words",
			sampleBits:  10,
			codedBits:   8,
			expected: Config{
				SampleStates:       1024,
				CodedStates:        256,
				LossyCodeWidth:     7,
				LossyRounding:      3,
				LosslessCodes:      128,
				LossyCodes:         128,
				CodeTables:         129,
				DefaultTable:       64,
				LowestTableThresh:  64,
				HighestTableThresh: 960,
				LossyFlag:          128,
			},
		},
	}

	for _, testCase := range tests {
		actual := SetupBMIP10(testCase.sampleBits, testCase.codedBits)

		if !reflect.DeepEqual(actual, &testCase.expected) {
			t.Errorf("Test '%s' failed to setup config properly", testCase.description)
			fmt.Printf("expected: %v\n", testCase.expected)
			fmt.Printf("actual: %v\n", actual)
		}
	}
}
