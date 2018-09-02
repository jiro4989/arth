package math

import (
	"bufio"
	"io"
	"math"
	"strconv"
	"strings"
)

// MinMaxSumAvg は入力から最小値、最大値、合計値、平均値を算出する
// needValuesフラグがtrueのときは入力をfloat64スライスに変換した値も返す
// needValuesフラグをセットしなければスライスは初期値のまま返却し、
// スライスにデータを保持しないため省メモリになる
func MinMaxSumAvg(r io.Reader, needValues bool) (cnt int, min, max, sum, avg float64, ns []float64, err error) {
	min = math.MaxFloat64 // 最初にでかい値を入れてないと判定されない
	max = 0.0
	sum = 0.0
	avg = 0.0

	// 入力をfloatに変換して都度計算
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		line = strings.Trim(line, " ")
		n, err := strconv.ParseFloat(line, 64)
		if err != nil {
			// 不正な文字列が存在しても後続の処理を継続してほしいのでcontinue
			continue
		}
		min = math.Min(n, min)
		max = math.Max(n, max)
		sum += n
		if needValues {
			ns = append(ns, n)
		}
		cnt++
	}
	if cnt == 0 {
		min = 0
		return
	}
	avg = sum / float64(cnt)
	err = sc.Err()
	return
}

// Median はfloat配列から中央値を算出する。
func Median(ns []float64) (med float64) {
	l := len(ns)
	if l <= 0 {
		return 0.0
	}
	if l%2 == 1 {
		return ns[l/2]
	}
	return ns[l/2-1]
}
