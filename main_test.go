package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/jiro4989/arth/internal/options"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	args := []string{"testdata/bigdata.txt"}
	opts := options.Options{
		CountFlag:         true,
		MinFlag:           true,
		MaxFlag:           true,
		SumFlag:           true,
		AverageFlag:       true,
		MedianFlag:        true,
		SordedFlag:        false,
		NoHeaderFlag:      false,
		OutFieldSeparator: "\t",
		OutFile:           "",
	}
	s, err := format(args, opts)
	assert.NoError(t, err)
	assert.Equal(t, "count\tmin\tmax\tsum\tavg\tmedian\n100\t1\t100\t5050\t50.5\t50", s)
}

func TestCalcOutValues(t *testing.T) {
	f := func(s string) io.Reader {
		return bytes.NewBufferString(s)
	}

	opts := options.Options{
		CountFlag:         true,
		MinFlag:           true,
		MaxFlag:           true,
		SumFlag:           true,
		AverageFlag:       true,
		MedianFlag:        true,
		SordedFlag:        false,
		NoHeaderFlag:      false,
		OutFieldSeparator: "\t",
		OutFile:           "",
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

func TestCalcMinMaxSumAvg(t *testing.T) {
	f := func(s string) io.Reader {
		return bytes.NewBufferString(s)
	}

	cnt, min, max, sum, avg, ns, err := calcMinMaxSumAvg(f("1.0\n2.0\n3.0\n4.0\n5.0\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)
	assert.Equal(t, 0, len(ns))

	cnt, min, max, sum, avg, ns, err = calcMinMaxSumAvg(f("5.0\n4.0\n3.0\n2.0\n1.0\n"), true)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)
	assert.EqualValues(t, []float64{5.0, 4.0, 3.0, 2.0, 1.0}, ns)

	cnt, min, max, sum, avg, ns, err = calcMinMaxSumAvg(f("4.0\n2.0\n3.0\n5.0\n1.0\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = calcMinMaxSumAvg(f("4.0\n2.0\n3.0\n5.0\n1.0\na\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = calcMinMaxSumAvg(f("foobar\n4.0\n2.0\n3.0\n5.0\n1.0\na\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = calcMinMaxSumAvg(f("1.0\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 1, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 1.0, max)
	assert.Equal(t, 1.0, sum)
	assert.Equal(t, 1.0, avg)

	cnt, min, max, sum, avg, ns, err = calcMinMaxSumAvg(f("1.0\n2.0\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 2, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 2.0, max)
	assert.Equal(t, 3.0, sum)
	assert.Equal(t, 1.5, avg)

	cnt, min, max, sum, avg, ns, err = calcMinMaxSumAvg(f("\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, 0.0, min)
	assert.Equal(t, 0.0, max)
	assert.Equal(t, 0.0, sum)
	assert.Equal(t, 0.0, avg)
}

func TestCalcMedian(t *testing.T) {
	assert.Equal(t, 3.0, calcMedian([]float64{1.0, 2.0, 3.0, 4.0, 5.0}))
	assert.Equal(t, 3.0, calcMedian([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}))
	assert.Equal(t, 5.0, calcMedian([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}))
	assert.Equal(t, 1.0, calcMedian([]float64{1.0}))
	assert.Equal(t, 1.0, calcMedian([]float64{1.0, 2.0}))
	assert.Equal(t, 0.0, calcMedian([]float64{}))
}
