package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/jiro4989/arth/internal/options"
	"github.com/stretchr/testify/assert"
)

func TestProcessInput(t *testing.T) {
	args := []string{"testdata/bigdata.txt"}
	opts := options.Options{
		CountFlag:    true,
		MinFlag:      true,
		MaxFlag:      true,
		SumFlag:      true,
		AverageFlag:  true,
		MedianFlag:   true,
		SordedFlag:   false,
		NoHeaderFlag: false,
		Delimiter:    "\t",
		OutFile:      "",
	}
	o, err := processInput(args, opts)
	assert.NoError(t, err)
	assert.Equal(t, []options.OutValues{
		options.OutValues{
			Count:   100,
			Min:     1,
			Max:     100,
			Sum:     5050,
			Average: 50.5,
			Median:  50,
		},
	}, o)

	os.Stdin, err = os.Open("testdata/normal_num.txt")
	assert.NoError(t, err)

	args = []string{}
	o, err = processInput(args, opts)
	assert.NoError(t, err)
	assert.Equal(t, []options.OutValues{
		options.OutValues{
			Count:   5,
			Min:     1,
			Max:     5,
			Sum:     15,
			Average: 3,
			Median:  3,
		},
	}, o)
}

func TestProcessStdin(t *testing.T) {
	opts := options.Options{
		CountFlag:    true,
		MinFlag:      true,
		MaxFlag:      true,
		SumFlag:      true,
		AverageFlag:  true,
		MedianFlag:   true,
		SordedFlag:   false,
		NoHeaderFlag: false,
		Delimiter:    "\t",
		OutFile:      "",
	}

	var err error
	os.Stdin, err = os.Open("testdata/normal_num.txt")
	assert.NoError(t, err)

	o, err := processStdin(opts)
	assert.NoError(t, err)
	assert.Equal(t, []options.OutValues{
		options.OutValues{
			Count:   5,
			Min:     1,
			Max:     5,
			Sum:     15,
			Average: 3,
			Median:  3,
		},
	}, o)
}

func TestProcessMultiInput(t *testing.T) {
	args := []string{
		"testdata/bigdata.txt",
		"testdata/normal_num.txt",
	}
	opts := options.Options{
		CountFlag:    true,
		MinFlag:      true,
		MaxFlag:      true,
		SumFlag:      true,
		AverageFlag:  true,
		MedianFlag:   true,
		SordedFlag:   false,
		NoHeaderFlag: false,
		Delimiter:    "\t",
		OutFile:      "",
	}
	o := processMultiInput(args, opts)
	assert.Equal(t, []options.OutValues{
		options.OutValues{
			Count:   100,
			Min:     1,
			Max:     100,
			Sum:     5050,
			Average: 50.5,
			Median:  50,
		},
		options.OutValues{
			Count:   5,
			Min:     1,
			Max:     5,
			Sum:     15,
			Average: 3,
			Median:  3,
		},
	}, o)
}

func TestCalcOutValues(t *testing.T) {
	f := func(s string) io.Reader {
		return bytes.NewBufferString(s)
	}

	opts := options.Options{
		CountFlag:    true,
		MinFlag:      true,
		MaxFlag:      true,
		SumFlag:      true,
		AverageFlag:  true,
		MedianFlag:   true,
		SordedFlag:   false,
		NoHeaderFlag: false,
		Delimiter:    "\t",
		OutFile:      "",
	}
	ov, err := calcOutValues(f("1.0\n2.0\n3.0\n4.0\n5.0\n"), opts)
	assert.NoError(t, err)
	assert.Equal(t, 5, ov.Count)
	assert.Equal(t, 1.0, ov.Min)
	assert.Equal(t, 5.0, ov.Max)
	assert.Equal(t, 15.0, ov.Sum)
	assert.Equal(t, 3.0, ov.Average)
	assert.Equal(t, 3.0, ov.Median)
}

func TestOut(t *testing.T) {
	lines := []string{
		"1",
		"2",
		"3",
	}
	opts := options.Options{
		Delimiter: "\t",
	}
	assert.NoError(t, out(lines, opts))

	opts = options.Options{
		Delimiter: "\t",
		OutFile:   "testdata/out/out.tsv",
	}
	assert.NoError(t, out(lines, opts))
}
