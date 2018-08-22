package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/jiro4989/arth/internal/options"
)

var Version string

func main() {
	var opts options.Options
	opts.Version = func() {
		fmt.Println(Version)
		os.Exit(0)
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}
	opts = opts.Setup()

	s, err := format(args, opts)
	if err != nil {
		panic(err)
	}

	if opts.OutFile == "" {
		fmt.Println(s)
	} else {
		err := ioutil.WriteFile(opts.OutFile, []byte(s), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

// format は引数、オプションを判定して計算し、出力する文字列を生成する。
func format(args []string, opts options.Options) (string, error) {
	// 引数指定がある場合はファイル名としてファイル読み込みを実施
	// 指定がない場合は標準入力を受け取る
	var r *os.File
	if 1 <= len(args) {
		var err error
		r, err = os.Open(args[0])
		if err != nil {
			return "", err
		}
		defer r.Close()
	} else {
		r = os.Stdin
	}

	ov, err := calcOutValues(r, opts)
	if err != nil {
		return "", err
	}
	return opts.Format(ov), nil
}

// calcOutValues は入力から出力データを計算する。
// オプションMedianFlagが存在するとき、ソートとソートデータの保持により
// メモリ消費と計算時間が増加する。
// オプションSordedFlagが存在するとき、入力がすでにソート済みとして
// ソート処理をスキップする。
func calcOutValues(r io.Reader, opts options.Options) (options.OutValues, error) {
	ov := options.OutValues{}
	ns := make([]float64, 0)
	var err error
	ov.Count, ov.Min, ov.Max, ov.Sum, ov.Average, ns, err = calcMinMaxSumAvg(r, opts.MedianFlag)
	if err != nil {
		return ov, err
	}

	if opts.MedianFlag {
		// SordedFlagがなければ、ソートを実行
		// SordedFlagがあれば、ソート済みとしてソートはスキップ(高速化)
		if !opts.SordedFlag {
			sort.Float64s(ns)
		}
		ov.Median = calcMedian(ns)
	}
	return ov, nil
}

// calcMinMaxSumAvg は入力から最小値、最大値、合計値、平均値を算出する
// needValuesフラグがtrueのときは入力をfloat64スライスに変換した値も返す
// needValuesフラグをセットしなければスライスは初期値のまま返却し、
// スライスにデータを保持しないため省メモリになる
func calcMinMaxSumAvg(r io.Reader, needValues bool) (cnt int, min, max, sum, avg float64, ns []float64, err error) {
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

// calcMedian はfloat配列から中央値を算出する。
func calcMedian(ns []float64) (med float64) {
	l := len(ns)
	if l <= 0 {
		return 0.0
	}
	if l%2 == 1 {
		return ns[l/2]
	}
	return ns[l/2-1]
}
