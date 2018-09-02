package options

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testdata struct {
	in     Options
	expect Options
}

type TestParseData struct {
	args    []string
	outopts Options
	outargs []string
}

func TestParse(t *testing.T) {
	tds := []TestParseData{
		TestParseData{
			args: []string{
				"main.go",
			},
			outopts: Options{
				CountFlag:    true,
				MinFlag:      true,
				MaxFlag:      true,
				SumFlag:      true,
				AverageFlag:  true,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
			outargs: []string{},
		},
		TestParseData{
			args: []string{
				"main.go",
				"--count",
			},
			outopts: Options{
				CountFlag:    true,
				MinFlag:      false,
				MaxFlag:      false,
				SumFlag:      false,
				AverageFlag:  false,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
			outargs: []string{},
		},
		TestParseData{
			args: []string{
				"main.go",
				"--min",
				"testdata",
			},
			outopts: Options{
				CountFlag:    false,
				MinFlag:      true,
				MaxFlag:      false,
				SumFlag:      false,
				AverageFlag:  false,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
			outargs: []string{"testdata"},
		},
	}
	for _, v := range tds {
		os.Args = v.args
		opts, args := Parse("1.0.0")
		assert.Equal(t, v.outopts.CountFlag, opts.CountFlag)
		assert.Equal(t, v.outopts.MinFlag, opts.MinFlag)
		assert.Equal(t, v.outopts.MaxFlag, opts.MaxFlag)
		assert.Equal(t, v.outopts.SumFlag, opts.SumFlag)
		assert.Equal(t, v.outopts.AverageFlag, opts.AverageFlag)
		assert.Equal(t, v.outopts.MedianFlag, opts.MedianFlag)
		assert.Equal(t, v.outopts.SordedFlag, opts.SordedFlag)
		assert.Equal(t, v.outopts.NoHeaderFlag, opts.NoHeaderFlag)
		assert.Equal(t, v.outopts.Delimiter, opts.Delimiter)
		assert.Equal(t, v.outargs, args)
	}
}

func TestSetup(t *testing.T) {
	tds := []testdata{
		testdata{
			in: Options{
				CountFlag:    false,
				MinFlag:      false,
				MaxFlag:      false,
				SumFlag:      false,
				AverageFlag:  false,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
			expect: Options{
				CountFlag:    true,
				MinFlag:      true,
				MaxFlag:      true,
				SumFlag:      true,
				AverageFlag:  true,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
		},
		testdata{
			in: Options{
				CountFlag:    true,
				MinFlag:      false,
				MaxFlag:      false,
				SumFlag:      false,
				AverageFlag:  false,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
			expect: Options{
				CountFlag:    true,
				MinFlag:      false,
				MaxFlag:      false,
				SumFlag:      false,
				AverageFlag:  false,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
		},
		testdata{
			in: Options{
				CountFlag:    false,
				MinFlag:      true,
				MaxFlag:      false,
				SumFlag:      true,
				AverageFlag:  false,
				MedianFlag:   false,
				SordedFlag:   true,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
			expect: Options{
				CountFlag:    false,
				MinFlag:      true,
				MaxFlag:      false,
				SumFlag:      true,
				AverageFlag:  false,
				MedianFlag:   false,
				SordedFlag:   true,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
		},
		testdata{
			in: Options{
				CountFlag:    true,
				MinFlag:      true,
				MaxFlag:      true,
				SumFlag:      true,
				AverageFlag:  true,
				MedianFlag:   false,
				SordedFlag:   true,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
			expect: Options{
				CountFlag:    true,
				MinFlag:      true,
				MaxFlag:      true,
				SumFlag:      true,
				AverageFlag:  true,
				MedianFlag:   false,
				SordedFlag:   true,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
		},
	}
	for _, v := range tds {
		v.in.Setup()
		assert.Equal(t, v.expect.CountFlag, v.in.CountFlag)
		assert.Equal(t, v.expect.MinFlag, v.in.MinFlag)
		assert.Equal(t, v.expect.MaxFlag, v.in.MaxFlag)
		assert.Equal(t, v.expect.SumFlag, v.in.SumFlag)
		assert.Equal(t, v.expect.AverageFlag, v.in.AverageFlag)
		assert.Equal(t, v.expect.MedianFlag, v.in.MedianFlag)
		assert.Equal(t, v.expect.SordedFlag, v.in.SordedFlag)
		assert.Equal(t, v.expect.NoHeaderFlag, v.in.NoHeaderFlag)
		assert.Equal(t, v.expect.Delimiter, v.in.Delimiter)
	}
}

type TestFormatData struct {
	ovs  []OutValues
	opts Options
	out  []string
}

func TestFormat(t *testing.T) {
	tds := []TestFormatData{
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
				OutValues{
					Count:   100,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:    true,
				MinFlag:      true,
				MaxFlag:      true,
				SumFlag:      true,
				AverageFlag:  true,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
			out: []string{
				"count\tmin\tmax\tsum\tavg",
				"2\t1\t2\t3\t1.5",
				"100\t1\t2\t3\t1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:    true,
				MinFlag:      true,
				MaxFlag:      true,
				SumFlag:      true,
				AverageFlag:  true,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    "\t",
			},
			out: []string{
				"count\tmin\tmax\tsum\tavg",
				"2\t1\t2\t3\t1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:    true,
				MinFlag:      true,
				MaxFlag:      true,
				SumFlag:      true,
				AverageFlag:  true,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: false,
				Delimiter:    ",",
			},
			out: []string{
				"count,min,max,sum,avg",
				"2,1,2,3,1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:    true,
				MinFlag:      true,
				MaxFlag:      true,
				SumFlag:      true,
				AverageFlag:  true,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: true,
				Delimiter:    ",",
			},
			out: []string{
				"2,1,2,3,1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:    false,
				MinFlag:      true,
				MaxFlag:      false,
				SumFlag:      true,
				AverageFlag:  true,
				MedianFlag:   false,
				SordedFlag:   false,
				NoHeaderFlag: true,
				Delimiter:    ",",
			},
			out: []string{
				"1,3,1.5",
			},
		},
	}

	for _, v := range tds {
		assert.Equal(t, v.out, Format(v.ovs, v.opts))
	}
}
