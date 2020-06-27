package typer

import (
	"reflect"
	"testing"
)

type DivideTextTestCase struct {
	text        string
	lines       []string
	framesCount int
}

func TestDivideTextOnLines(t *testing.T) {
	const (
		W             = 500
		H             = 250
		maxLineSize   = 5
		maxLinesCount = 5
	)

	testCases := []DivideTextTestCase{
		{"This is a test line",
			[]string{"This ", "is a ", "test ", "line"},
			19,
		},
	}

	generator := &Typer{
		maxLinesCount: maxLinesCount,
		maxLineSize:   maxLineSize,
		frameW:        W,
		frameH:        H,
	}

	for _, testCase := range testCases {
		lines, framesCount, err := generator.divideTextOnLines(testCase.text)
		if err != nil {
			t.Error(err.Error())
		}
		if framesCount != testCase.framesCount {
			t.Errorf("number of frames is not right: %v != %v", framesCount, testCase.framesCount)
		}
		if !reflect.DeepEqual(lines, testCase.lines) {
			t.Errorf("failed: %v != %v", lines, testCase.lines)
		}
	}
}
