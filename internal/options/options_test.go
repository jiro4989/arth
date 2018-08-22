package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testdata struct {
	in     Options
	expect Options
}

func TestSetup(t *testing.T) {
	tds := []testdata{
		testdata{
			in: Options{
				CountFlag:         false,
				MinFlag:           false,
				MaxFlag:           false,
				SumFlag:           false,
				AverageFlag:       false,
				MedianFlag:        false,
				SordedFlag:        false,
				NoHeaderFlag:      false,
				OutFieldSeparator: "\t",
			},
			expect: Options{
				CountFlag:         true,
				MinFlag:           true,
				MaxFlag:           true,
				SumFlag:           true,
				AverageFlag:       true,
				MedianFlag:        false,
				SordedFlag:        false,
				NoHeaderFlag:      false,
				OutFieldSeparator: "\t",
			},
		},
		testdata{
			in: Options{
				CountFlag:         true,
				MinFlag:           false,
				MaxFlag:           false,
				SumFlag:           false,
				AverageFlag:       false,
				MedianFlag:        false,
				SordedFlag:        false,
				NoHeaderFlag:      false,
				OutFieldSeparator: "\t",
			},
			expect: Options{
				CountFlag:         true,
				MinFlag:           false,
				MaxFlag:           false,
				SumFlag:           false,
				AverageFlag:       false,
				MedianFlag:        false,
				SordedFlag:        false,
				NoHeaderFlag:      false,
				OutFieldSeparator: "\t",
			},
		},
		testdata{
			in: Options{
				CountFlag:         false,
				MinFlag:           true,
				MaxFlag:           false,
				SumFlag:           true,
				AverageFlag:       false,
				MedianFlag:        false,
				SordedFlag:        true,
				NoHeaderFlag:      false,
				OutFieldSeparator: "\t",
			},
			expect: Options{
				CountFlag:         false,
				MinFlag:           true,
				MaxFlag:           false,
				SumFlag:           true,
				AverageFlag:       false,
				MedianFlag:        false,
				SordedFlag:        true,
				NoHeaderFlag:      false,
				OutFieldSeparator: "\t",
			},
		},
		testdata{
			in: Options{
				CountFlag:         true,
				MinFlag:           true,
				MaxFlag:           true,
				SumFlag:           true,
				AverageFlag:       true,
				MedianFlag:        false,
				SordedFlag:        true,
				NoHeaderFlag:      false,
				OutFieldSeparator: "\t",
			},
			expect: Options{
				CountFlag:         true,
				MinFlag:           true,
				MaxFlag:           true,
				SumFlag:           true,
				AverageFlag:       true,
				MedianFlag:        false,
				SordedFlag:        true,
				NoHeaderFlag:      false,
				OutFieldSeparator: "\t",
			},
		},
	}
	for _, v := range tds {
		got := v.in.Setup()
		assert.Equal(t, v.expect.CountFlag, got.CountFlag)
		assert.Equal(t, v.expect.MinFlag, got.MinFlag)
		assert.Equal(t, v.expect.MaxFlag, got.MaxFlag)
		assert.Equal(t, v.expect.SumFlag, got.SumFlag)
		assert.Equal(t, v.expect.AverageFlag, got.AverageFlag)
		assert.Equal(t, v.expect.MedianFlag, got.MedianFlag)
		assert.Equal(t, v.expect.SordedFlag, got.SordedFlag)
		assert.Equal(t, v.expect.NoHeaderFlag, got.NoHeaderFlag)
		assert.Equal(t, v.expect.OutFieldSeparator, got.OutFieldSeparator)
	}
}

func TestFormat(t *testing.T) {
	v := OutValues{
		Count:   2,
		Min:     1.0,
		Max:     2.0,
		Sum:     3.0,
		Average: 1.5,
		Median:  1.0,
	}
	opts := Options{
		CountFlag:         true,
		MinFlag:           true,
		MaxFlag:           true,
		SumFlag:           true,
		AverageFlag:       true,
		MedianFlag:        false,
		SordedFlag:        false,
		NoHeaderFlag:      false,
		OutFieldSeparator: "\t",
	}

	assert.Equal(t, "count\tmin\tmax\tsum\tavg\n2\t1\t2\t3\t1.5", opts.Format(v))
	opts.OutFieldSeparator = ","
	assert.Equal(t, "count,min,max,sum,avg\n2,1,2,3,1.5", opts.Format(v))
}
