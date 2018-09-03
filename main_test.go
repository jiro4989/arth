package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/jiro4989/arth/internal/options"
	"github.com/stretchr/testify/assert"
)

type TestMainData struct {
	args []string
}

func TestMain(t *testing.T) {
	tds := []TestMainData{
		TestMainData{
			args: []string{
				"main.go",
				"testdata/normal_num.txt",
			},
		},
		TestMainData{
			args: []string{
				"main.go",
				"testdata/normal_num.txt",
				"testdata/normal_num.txt",
			},
		},
		TestMainData{
			args: []string{
				"main.go",
				"-D", ",",
				"testdata/normal_num.txt",
			},
		},
		TestMainData{
			args: []string{
				"main.go",
				"-m",
				"testdata/bigdata.txt",
			},
		},
		TestMainData{
			args: []string{
				"main.go",
				"-p",
				"95",
				"testdata/bigdata.txt",
			},
		},
		TestMainData{
			args: []string{
				"main.go",
				"-f", "1:testdata/sample.csv",
				"-f", "2:testdata/sample.csv",
				"-f", "testdata/sample.csv",
			},
		},
	}
	for _, v := range tds {
		os.Args = v.args
		main()
	}
}

func TestProcessInput(t *testing.T) {
	args := []string{"testdata/bigdata.txt"}
	opts := options.Options{
		CountFlag:      true,
		MinFlag:        true,
		MaxFlag:        true,
		SumFlag:        true,
		AverageFlag:    true,
		MedianFlag:     true,
		SortedFlag:     false,
		HeaderFlag:     false,
		InputDelimiter: "\t",
		OutFile:        "",
	}
	o, err := processInput(args, opts)
	assert.NoError(t, err)
	assert.Equal(t, []options.OutValues{
		options.OutValues{
			FileName: "testdata/bigdata.txt",
			Count:    100,
			Min:      1,
			Max:      100,
			Sum:      5050,
			Average:  50.5,
			Median:   50,
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
		CountFlag:      true,
		MinFlag:        true,
		MaxFlag:        true,
		SumFlag:        true,
		AverageFlag:    true,
		MedianFlag:     true,
		SortedFlag:     false,
		HeaderFlag:     false,
		InputDelimiter: "\t",
		OutFile:        "",
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

type TestProcessMultiInputData struct {
	args []string
	opts options.Options
	out  []options.OutValues
}

func TestProcessMultiInput(t *testing.T) {
	tds := []TestProcessMultiInputData{
		TestProcessMultiInputData{
			args: []string{
				"testdata/bigdata.txt",
				"testdata/normal_num.txt",
			},
			opts: options.Options{
				CountFlag:      true,
				MinFlag:        true,
				MaxFlag:        true,
				SumFlag:        true,
				AverageFlag:    true,
				MedianFlag:     true,
				SortedFlag:     false,
				HeaderFlag:     false,
				InputDelimiter: "\t",
				OutFile:        "",
			},
			out: []options.OutValues{
				options.OutValues{
					FileName: "testdata/bigdata.txt",
					Count:    100,
					Min:      1,
					Max:      100,
					Sum:      5050,
					Average:  50.5,
					Median:   50,
				},
				options.OutValues{
					FileName: "testdata/normal_num.txt",
					Count:    5,
					Min:      1,
					Max:      5,
					Sum:      15,
					Average:  3,
					Median:   3,
				},
			},
		},
		TestProcessMultiInputData{
			args: []string{
				"testdata/bigdata.txt",
				"testdata/normal_num.txt",
			},
			opts: options.Options{
				CountFlag:      true,
				MinFlag:        true,
				MaxFlag:        true,
				SumFlag:        true,
				AverageFlag:    true,
				MedianFlag:     true,
				SortedFlag:     false,
				HeaderFlag:     false,
				Percentile:     95,
				InputDelimiter: "\t",
				OutFile:        "",
			},
			out: []options.OutValues{
				options.OutValues{
					FileName:   "testdata/bigdata.txt",
					Count:      100,
					Min:        1,
					Max:        100,
					Sum:        5050,
					Average:    50.5,
					Median:     50,
					Percentile: 95,
				},
				options.OutValues{
					FileName:   "testdata/normal_num.txt",
					Count:      5,
					Min:        1,
					Max:        5,
					Sum:        15,
					Average:    3,
					Median:     3,
					Percentile: 4,
				},
			},
		},
		TestProcessMultiInputData{
			args: []string{
				"testdata/bigdata.txt",
				"testdata/normal_num.txt",
			},
			opts: options.Options{
				CountFlag:      true,
				MinFlag:        true,
				MaxFlag:        true,
				SumFlag:        true,
				AverageFlag:    true,
				MedianFlag:     false,
				SortedFlag:     false,
				HeaderFlag:     false,
				Percentile:     95,
				InputDelimiter: "\t",
				OutFile:        "",
			},
			out: []options.OutValues{
				options.OutValues{
					FileName:   "testdata/bigdata.txt",
					Count:      100,
					Min:        1,
					Max:        100,
					Sum:        5050,
					Average:    50.5,
					Percentile: 95,
				},
				options.OutValues{
					FileName:   "testdata/normal_num.txt",
					Count:      5,
					Min:        1,
					Max:        5,
					Sum:        15,
					Average:    3,
					Percentile: 4,
				},
			},
		},
	}
	for _, v := range tds {
		o := processMultiInput(v.args, v.opts)
		assert.Equal(t, v.out, o)
	}
}

type TestCalcOutValuesData struct {
	r    io.Reader
	opts options.Options
	out  options.OutValues
}

func TestCalcOutValues(t *testing.T) {
	f := func(ss ...string) io.Reader {
		return bytes.NewBufferString(strings.Join(ss, "\n"))
	}

	tds := []TestCalcOutValuesData{
		TestCalcOutValuesData{
			r: f(
				"1.0",
				"2.0",
				"3.0",
				"4.0",
				"5.0",
			),
			opts: options.Options{
				CountFlag:   true,
				MinFlag:     true,
				MaxFlag:     true,
				SumFlag:     true,
				AverageFlag: true,
				MedianFlag:  true,
				Percentile:  95,
			},
			out: options.OutValues{
				Count:      5,
				Min:        1,
				Max:        5,
				Sum:        15,
				Average:    3,
				Median:     3,
				Percentile: 4,
			},
		},
		TestCalcOutValuesData{
			r: f(
				"1.0",
				"2.0",
				"3.0",
				"4.0",
				"5.0",
			),
			opts: options.Options{
				CountFlag:   true,
				MinFlag:     true,
				MaxFlag:     true,
				SumFlag:     true,
				AverageFlag: true,
				MedianFlag:  true,
			},
			out: options.OutValues{
				Count:   5,
				Min:     1,
				Max:     5,
				Sum:     15,
				Average: 3,
				Median:  3,
			},
		},
		TestCalcOutValuesData{
			r: f(
				"1.0",
				"2.0",
				"3.0",
				"4.0",
				"5.0",
			),
			opts: options.Options{
				CountFlag:  true,
				Percentile: 95,
			},
			out: options.OutValues{
				Count:      5,
				Min:        1,
				Max:        5,
				Sum:        15,
				Average:    3,
				Percentile: 4,
			},
		},
		TestCalcOutValuesData{
			r: f(
				"1.0",
				"2.0",
				"3.0",
				"4.0",
				"5.0",
			),
			opts: options.Options{
				MinFlag: true,
			},
			out: options.OutValues{
				Count:   5,
				Min:     1,
				Max:     5,
				Sum:     15,
				Average: 3,
			},
		},
		TestCalcOutValuesData{
			r: f(
				"1.0",
				"2.0",
				"5.0",
				"3.0",
				"4.0",
			),
			opts: options.Options{
				SortedFlag: true,
				MedianFlag: true,
				Percentile: 95,
			},
			out: options.OutValues{
				Count:      5,
				Min:        1,
				Max:        5,
				Sum:        15,
				Average:    3,
				Median:     5,
				Percentile: 3,
			},
		},
		TestCalcOutValuesData{
			r: f(
				"1.0",
				"2.0",
				"5.0",
				"3.0",
				"4.0",
			),
			opts: options.Options{
				SortedFlag: true,
				Percentile: 95,
			},
			out: options.OutValues{
				Count:      5,
				Min:        1,
				Max:        5,
				Sum:        15,
				Average:    3,
				Percentile: 3,
			},
		},
		TestCalcOutValuesData{
			r: f(
				"1.0",
				"2.0",
				"5.0",
				"3.0",
				"4.0",
			),
			opts: options.Options{
				SortedFlag: true,
				MedianFlag: true,
			},
			out: options.OutValues{
				Count:   5,
				Min:     1,
				Max:     5,
				Sum:     15,
				Average: 3,
				Median:  5,
			},
		},
	}

	for _, v := range tds {
		ov, err := calcOutValues(v.r, v.opts, nil)
		assert.NoError(t, err)
		assert.Equal(t, v.out, ov)
	}
}

func TestOut(t *testing.T) {
	lines := []string{
		"1",
		"2",
		"3",
	}
	opts := options.Options{
		InputDelimiter: "\t",
	}
	assert.NoError(t, out(lines, opts))

	opts = options.Options{
		InputDelimiter: "\t",
		OutFile:        "testdata/out/out.tsv",
	}
	assert.NoError(t, out(lines, opts))
}
