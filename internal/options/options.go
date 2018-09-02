package options

import (
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
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

// Parse はコマンドラインオプションを解析する。
// 解析あとはオプションと、残った引数を返す。
// また、入力ファイルパスの重複を除外する。
func Parse(version string) (Options, []string) {
	var opts Options
	opts.Version = func() {
		fmt.Println(version)
		os.Exit(0)
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}
	opts.Setup()

	return opts, args
}

// Setup はオプションのデフォルト値をセットします。
// Count, Min, Max, Sumのいずれもfalseの場合は、すべてtrueにする。
func (o *Options) Setup() {
	if !o.CountFlag && !o.MinFlag && !o.MaxFlag && !o.SumFlag && !o.AverageFlag {
		o.CountFlag = true
		o.MinFlag = true
		o.SumFlag = true
		o.MaxFlag = true
		o.AverageFlag = true
	}
}

// Format は出力用のデータをオプションに応じて出力ように整形する。
func Format(vs []OutValues, opts Options) []string {
	maps := make([]map[string]string, len(vs))
	for i, v := range vs {
		m := make(map[string]string) // 値
		setFunc := func(flg bool, h string, v interface{}) {
			if flg {
				if n, ok := v.(int); ok {
					s := fmt.Sprintf("%d", n)
					m[h] = s
					return
				}

				if n, ok := v.(float64); ok {
					s := fmt.Sprintf("%.2f", n)
					// 不要な末尾の0埋めを削除
					s = strings.TrimRight(s, "0")
					s = strings.TrimRight(s, ".")
					m[h] = s
					return
				}

				s := fmt.Sprintf("%v", v)
				m[h] = s
			}
		}

		setFunc(opts.CountFlag, HeaderCount, v.Count)
		setFunc(opts.MinFlag, HeaderMin, v.Min)
		setFunc(opts.MaxFlag, HeaderMax, v.Max)
		setFunc(opts.SumFlag, HeaderSum, v.Sum)
		setFunc(opts.AverageFlag, HeaderAverage, v.Average)
		setFunc(opts.MedianFlag, HeaderMedian, v.Median)

		maps[i] = m
	}

	lines := make([]string, 0)

	headers := make([]string, 0)
	// ヘッダの連結
	// オプションがあるとセットしない
	m := maps[0]
	for _, k := range []string{
		HeaderCount,
		HeaderMin,
		HeaderMax,
		HeaderSum,
		HeaderAverage,
		HeaderMedian,
	} {
		if m[k] != "" {
			headers = append(headers, k)
		}
	}

	s := strings.Join(headers, opts.Delimiter)
	if !opts.NoHeaderFlag {
		lines = append(lines, s)
	}

	// 値の追加
	for _, m := range maps {
		values := make([]string, 0)
		for _, k := range headers {
			values = append(values, m[k])
		}
		s := strings.Join(values, opts.Delimiter)
		lines = append(lines, s)
	}

	return lines
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
