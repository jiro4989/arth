package math

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinMaxSumAvg(t *testing.T) {
	f := func(s string) io.Reader {
		return bytes.NewBufferString(s)
	}

	cnt, min, max, sum, avg, ns, err := MinMaxSumAvg(f("1.0\n2.0\n3.0\n4.0\n5.0\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)
	assert.Equal(t, 0, len(ns))

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("5.0\n4.0\n3.0\n2.0\n1.0\n"), true)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)
	assert.EqualValues(t, []float64{5.0, 4.0, 3.0, 2.0, 1.0}, ns)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("4.0\n2.0\n3.0\n5.0\n1.0\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("4.0\n2.0\n3.0\n5.0\n1.0\na\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("foobar\n4.0\n2.0\n3.0\n5.0\n1.0\na\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("1.0\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 1, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 1.0, max)
	assert.Equal(t, 1.0, sum)
	assert.Equal(t, 1.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("1.0\n2.0\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 2, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 2.0, max)
	assert.Equal(t, 3.0, sum)
	assert.Equal(t, 1.5, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("\n"), false)
	assert.NoError(t, err)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, 0.0, min)
	assert.Equal(t, 0.0, max)
	assert.Equal(t, 0.0, sum)
	assert.Equal(t, 0.0, avg)
}

func TestMedian(t *testing.T) {
	assert.Equal(t, 3.0, Median([]float64{1.0, 2.0, 3.0, 4.0, 5.0}))
	assert.Equal(t, 3.0, Median([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}))
	assert.Equal(t, 5.0, Median([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}))
	assert.Equal(t, 1.0, Median([]float64{1.0}))
	assert.Equal(t, 1.0, Median([]float64{1.0, 2.0}))
	assert.Equal(t, 0.0, Median([]float64{}))
}
