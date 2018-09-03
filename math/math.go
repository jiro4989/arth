package math

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// MinMaxSumAvgConfig はMinMaxSumAvg関数の設定です。
type MinMaxSumAvgConfig struct {
	// NeedValues は計算途中に読み込んだデータを返却するか否かです。
	NeedValues bool
	// Delimiter は読み込んだ行データの区切り文字です。
	Delimiter string
	// FieldNum は読み込んだ行データのうち、取り出すフィールド番号です。
	FieldIndex int
	// IgnoreHeaderRows は読み込むデータの開始から無視する行数です。
	IgnoreHeaderRows int
}

// MinMaxSumAvg は入力から最小値、最大値、合計値、平均値を算出する
// needValuesフラグがtrueのときは入力をfloat64スライスに変換した値も返す
// needValuesフラグをセットしなければスライスは初期値のまま返却し、
// スライスにデータを保持しないため省メモリになる
func MinMaxSumAvg(r io.Reader, conf MinMaxSumAvgConfig) (cnt int, min, max, sum, avg float64, ns []float64, err error) {
	min = math.MaxFloat64 // 最初にでかい値を入れてないと判定されない
	max = 0.0
	sum = 0.0
	avg = 0.0

	ignoredCounter := 0
	// 入力をfloatに変換して都度計算
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		// 指定行数まで無視
		if ignoredCounter < conf.IgnoreHeaderRows {
			ignoredCounter++
			continue
		}

		line := sc.Text()
		line = strings.Trim(line, " ")
		line = cutField(line, conf.Delimiter, conf.FieldIndex)
		n, err := strconv.ParseFloat(line, 64)
		if err != nil {
			// 不正な文字列が存在しても後続の処理を継続してほしいのでcontinue
			msg := fmt.Sprintf("warn: illegal value. value=%v", line)
			fmt.Fprintln(os.Stderr, msg)
			continue
		}
		min = math.Min(n, min)
		max = math.Max(n, max)
		sum += n
		if conf.NeedValues {
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

// cutField は文字列を指定文字で区切り、指定の番号のフィールドを返す。
func cutField(l, d string, i int) string {
	if i <= 0 || l == "" {
		return l
	}

	ss := strings.Split(l, d)
	if len(ss) <= 1 {
		return l
	}
	n := i - 1
	if len(ss) <= n {
		return l
	}
	if n < 0 {
		n = 0
	}
	// 値指定は1〜なので-1する
	return ss[n]
}

// Median はfloat配列から中央値を算出する。
func Median(ns []float64) float64 {
	l := len(ns)
	if l <= 0 {
		return 0.0
	}
	if l%2 == 1 {
		return ns[l/2]
	}
	return ns[l/2-1]
}

// Percentile はパーセンタイル値を計算する。
func Percentile(ns []float64, n int) float64 {
	if n <= 0 {
		return 0.0
	}

	l := len(ns)
	if l <= 0 {
		return 0.0
	}

	i := l*n/100 - 1
	if i < 0 {
		i = 0
	}
	return ns[i]
}
