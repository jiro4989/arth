package math

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestMinMaxSumAvgData struct {
	inR          io.Reader
	inNeedValues bool
	inF          func(string) string
	outCount     int
	outMin       float64
	outMax       float64
	outSum       float64
	outAvg       float64
	outNs        []float64
}

func TestMinMaxSumAvg(t *testing.T) {
	f := func(s ...string) io.Reader {
		return bytes.NewBufferString(strings.Join(s, "\n"))
	}

	tds := []TestMinMaxSumAvgData{
		TestMinMaxSumAvgData{ // ソート不要のとき、nsが空になる
			inR: f(
				"1.0",
				"2.0",
				"3.0",
				"4.0",
				"5.0",
			),
			inNeedValues: false,
			inF:          nil,
			outCount:     5,
			outMin:       1.0,
			outMax:       5.0,
			outSum:       15.0,
			outAvg:       3.0,
			outNs:        nil,
		},
		TestMinMaxSumAvgData{ // ソートされてないデータとソート必要フラグ。読み取ったデータを返却
			inR: f(
				"1.0",
				"5.0",
				"4.0",
				"3.0",
				"2.0",
			),
			inNeedValues: true,
			inF:          nil,
			outCount:     5,
			outMin:       1.0,
			outMax:       5.0,
			outSum:       15.0,
			outAvg:       3.0,
			outNs:        []float64{1, 5, 4, 3, 2},
		},
		TestMinMaxSumAvgData{ // 不正なデータがあってもエラーを無視すること
			inR: f(
				"1.0",
				"5.0",
				"4.0",
				"a",
				"3.0",
				"2.0",
			),
			inNeedValues: true,
			inF:          nil,
			outCount:     5,
			outMin:       1.0,
			outMax:       5.0,
			outSum:       15.0,
			outAvg:       3.0,
			outNs:        []float64{1, 5, 4, 3, 2},
		},
		TestMinMaxSumAvgData{ // データが1つだけ
			inR: f(
				"1.0",
			),
			inNeedValues: true,
			inF:          nil,
			outCount:     1,
			outMin:       1.0,
			outMax:       1.0,
			outSum:       1.0,
			outAvg:       1.0,
			outNs:        []float64{1.0},
		},
		TestMinMaxSumAvgData{ // データが0
			inR:          f(),
			inNeedValues: true,
			inF:          nil,
			outCount:     0,
			outMin:       0.0,
			outMax:       0.0,
			outSum:       0.0,
			outAvg:       0.0,
			outNs:        nil,
		},
		TestMinMaxSumAvgData{ // 空データのみ
			inR: f(
				"",
				"",
			),
			inNeedValues: true,
			inF:          nil,
			outCount:     0,
			outMin:       0.0,
			outMax:       0.0,
			outSum:       0.0,
			outAvg:       0.0,
			outNs:        nil,
		},
		TestMinMaxSumAvgData{ // 関数適用確認
			inR: f(
				"val1,val2",
				"1,2",
				"3,4",
				"5,6",
			),
			inNeedValues: true,
			inF: func(s string) string {
				return strings.Split(s, ",")[0]
			},
			outCount: 3,
			outMin:   1.0,
			outMax:   5.0,
			outSum:   9.0,
			outAvg:   3.0,
			outNs:    []float64{1, 3, 5},
		},
	}
	for _, v := range tds {
		cnt, min, max, sum, avg, ns, err := MinMaxSumAvg(v.inR, v.inNeedValues, v.inF)
		assert.NoError(t, err)
		assert.Equal(t, v.outCount, cnt)
		assert.Equal(t, v.outMin, min)
		assert.Equal(t, v.outMax, max)
		assert.Equal(t, v.outSum, sum)
		assert.Equal(t, v.outAvg, avg)
		assert.EqualValues(t, v.outNs, ns)
	}
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
