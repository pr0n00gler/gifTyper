package typer

import (
	"testing"
)

type DivideTextTestCase struct {
	text        string
	linesCount  int
	framesCount int
}

func TestDivideTextOnLines(t *testing.T) {
	generator, err := InitGenerator()
	if err != nil {
		t.Error(err.Error())
		return
	}

	testCases := []DivideTextTestCase{
		{"This is a test line",
			1,
			19,
		},
	}

	testCase := testCases[0]
	lines, framesCount, err := generator.divideTextOnLines(testCase.text)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if framesCount != testCase.framesCount {
		t.Errorf("number of frames is not right: %v != %v", framesCount, testCase.framesCount)
	}

	if testCase.linesCount != len(lines) {
		t.Errorf("failed: %v != %v", testCase.linesCount, len(lines))
	}
}
