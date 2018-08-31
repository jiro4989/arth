package options

import (
	"fmt"
	"strings"
)

const (
	HeaderCount   = "count"
	HeaderMin     = "min"
	HeaderMax     = "max"
	HeaderSum     = "sum"
	HeaderAverage = "avg"
	HeaderMedian  = "median"
)

// Options はコマンドラインオプション引数です。
type Options struct {
	Version      func() `short:"v" long:"version" description:"バージョン情報"`
	CountFlag    bool   `long:"count" description:"データ数を出力する"`
	MinFlag      bool   `long:"min" description:"最小値を出力する"`
	MaxFlag      bool   `long:"max" description:"最大値を出力する"`
	SumFlag      bool   `long:"sum" description:"合計を出力する"`
	AverageFlag  bool   `long:"avg" description:"平均値を出力する"`
	MedianFlag   bool   `short:"m" long:"median" description:"中央値を出力する"`
	SordedFlag   bool   `short:"s" long:"sorted" description:"入力元データがソート済みフラグ"`
	NoHeaderFlag bool   `short:"n" long:"noheader" description:"ヘッダを出力しない"`
	Delimiter    string `short:"d" long:"delimiter" description:"出力時の区切り文字を指定" default:"\t"`
	OutFile      string `short:"o" long:"outfile" description:"出力ファイルパス"`
}

// Setup はオプションのデフォルト値をセットします。
// Count, Min, Max, Sumのいずれもfalseの場合は、すべてtrueにする。
func (o Options) Setup() Options {
	if !o.CountFlag && !o.MinFlag && !o.MaxFlag && !o.SumFlag && !o.AverageFlag {
		o.CountFlag = true
		o.MinFlag = true
		o.SumFlag = true
		o.MaxFlag = true
		o.AverageFlag = true
	}
	return o
}

// Format はオプションフラグに応じて標準出力します。
func (o *Options) Format(v OutValues) string {
	hs := make([]string, 0) // ヘッダ
	vs := make([]string, 0) // 値
	setFunc := func(flg bool, h string, v interface{}) {
		if flg {
			if !o.NoHeaderFlag {
				hs = append(hs, h)
			}
			if n, ok := v.(int); ok {
				vs = append(vs, fmt.Sprintf("%d", n))
				return
			}
			if n, ok := v.(float64); ok {
				s := fmt.Sprintf("%.2f", n)
				// 不要な末尾の0埋めを削除
				s = strings.TrimRight(s, "0")
				s = strings.TrimRight(s, ".")
				vs = append(vs, s)
				return
			}
			vs = append(vs, fmt.Sprintf("%v", v))
		}
	}

	setFunc(o.CountFlag, HeaderCount, v.Count)
	setFunc(o.MinFlag, HeaderMin, v.Min)
	setFunc(o.MaxFlag, HeaderMax, v.Max)
	setFunc(o.SumFlag, HeaderSum, v.Sum)
	setFunc(o.AverageFlag, HeaderAverage, v.Average)
	setFunc(o.MedianFlag, HeaderMedian, v.Median)

	record := strings.Join(vs, o.Delimiter)
	if !o.NoHeaderFlag {
		h := strings.Join(hs, o.Delimiter)
		record = h + "\n" + record
	}
	return record
}

// OutValues はFormat関数で使用する値構造体です。
type OutValues struct {
	Count   int
	Min     float64
	Max     float64
	Sum     float64
	Average float64
	Median  float64
}
