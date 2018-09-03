package math

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinMaxSumAvg(t *testing.T) {
	f := func(s string) io.Reader {
		return bytes.NewBufferString(s)
	}

	cnt, min, max, sum, avg, ns, err := MinMaxSumAvg(f("1.0\n2.0\n3.0\n4.0\n5.0\n"), false, nil)
	assert.NoError(t, err)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)
	assert.Equal(t, 0, len(ns))

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("5.0\n4.0\n3.0\n2.0\n1.0\n"), true, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)
	assert.EqualValues(t, []float64{5.0, 4.0, 3.0, 2.0, 1.0}, ns)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("4.0\n2.0\n3.0\n5.0\n1.0\n"), false, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("4.0\n2.0\n3.0\n5.0\n1.0\na\n"), false, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("foobar\n4.0\n2.0\n3.0\n5.0\n1.0\na\n"), false, nil)
	assert.NoError(t, err)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 15.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("1.0\n"), false, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 1.0, max)
	assert.Equal(t, 1.0, sum)
	assert.Equal(t, 1.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("1.0\n2.0\n"), false, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 2.0, max)
	assert.Equal(t, 3.0, sum)
	assert.Equal(t, 1.5, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("\n"), false, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, 0.0, min)
	assert.Equal(t, 0.0, max)
	assert.Equal(t, 0.0, sum)
	assert.Equal(t, 0.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("val1,val2\n1,2\n3,4\n5,6\n"), false, func(s string) string {
		return strings.Split(s, ",")[0]
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, cnt)
	assert.Equal(t, 1.0, min)
	assert.Equal(t, 5.0, max)
	assert.Equal(t, 9.0, sum)
	assert.Equal(t, 3.0, avg)

	cnt, min, max, sum, avg, ns, err = MinMaxSumAvg(f("val1,val2\n1,2\n3,4\n5,6\n"), false, func(s string) string {
		return strings.Split(s, ",")[1]
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, cnt)
	assert.Equal(t, 2.0, min)
	assert.Equal(t, 6.0, max)
	assert.Equal(t, 12.0, sum)
	assert.Equal(t, 4.0, avg)
}

func TestMedian(t *testing.T) {
	assert.Equal(t, 3.0, Median([]float64{1.0, 2.0, 3.0, 4.0, 5.0}))
	assert.Equal(t, 3.0, Median([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}))
	assert.Equal(t, 5.0, Median([]float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}))
	assert.Equal(t, 1.0, Median([]float64{1.0}))
	assert.Equal(t, 1.0, Median([]float64{1.0, 2.0}))
	assert.Equal(t, 0.0, Median([]float64{}))
}

type TestPercentileData struct {
	ns []float64
	n  int
	p  float64 // percentile
}

func TestPercentile(t *testing.T) {
	tds := []TestPercentileData{
		TestPercentileData{
			ns: []float64{
				1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0,
				11.0, 12.0, 13.0, 14.0, 15.0, 16.0, 17.0, 18.0, 19.0, 20.0,
				21.0, 22.0, 23.0, 24.0, 25.0, 26.0, 27.0, 28.0, 29.0, 30.0,
				31.0, 32.0, 33.0, 34.0, 35.0, 36.0, 37.0, 38.0, 39.0, 40.0,
				41.0, 42.0, 43.0, 44.0, 45.0, 46.0, 47.0, 48.0, 49.0, 50.0,
				51.0, 52.0, 53.0, 54.0, 55.0, 56.0, 57.0, 58.0, 59.0, 60.0,
				61.0, 62.0, 63.0, 64.0, 65.0, 66.0, 67.0, 68.0, 69.0, 70.0,
				71.0, 72.0, 73.0, 74.0, 75.0, 76.0, 77.0, 78.0, 79.0, 80.0,
				81.0, 82.0, 83.0, 84.0, 85.0, 86.0, 87.0, 88.0, 89.0, 90.0,
				91.0, 92.0, 93.0, 94.0, 95.0, 96.0, 97.0, 92228.0, 11199.0, 999900.0,
			},
			n: 95,
			p: 95.0,
		},
		TestPercentileData{
			ns: []float64{1.0},
			n:  95,
			p:  1.0,
		},
		TestPercentileData{
			ns: []float64{1.0},
			n:  0,
			p:  0.0,
		},
		TestPercentileData{
			ns: []float64{1.0},
			n:  -1,
			p:  0.0,
		},
		TestPercentileData{
			ns: []float64{},
			n:  95,
			p:  0.0,
		},
	}
	for _, v := range tds {
		assert.Equal(t, v.p, Percentile(v.ns, v.n))
	}
}
