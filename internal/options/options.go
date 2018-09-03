package options

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

const (
	FileName         = "filename"
	HeaderCount      = "count"
	HeaderMin        = "min"
	HeaderMax        = "max"
	HeaderSum        = "sum"
	HeaderAverage    = "avg"
	HeaderMedian     = "median"
	HeaderPercentile = "percentile"
)

// Options はコマンドラインオプション引数です。
type Options struct {
	Version             func()                `short:"v" long:"version" description:"バージョン情報"`
	NoFileNameFlag      bool                  `short:"N" long:"nofilename" description:"入力元ファイル名を出力しない"`
	CountFlag           bool                  `short:"c" long:"count" description:"データ数を出力する"`
	MinFlag             bool                  `short:"n" long:"min" description:"最小値を出力する"`
	MaxFlag             bool                  `short:"x" long:"max" description:"最大値を出力する"`
	SumFlag             bool                  `short:"u" long:"sum" description:"合計を出力する"`
	AverageFlag         bool                  `short:"a" long:"avg" description:"平均値を出力する"`
	MedianFlag          bool                  `short:"m" long:"median" description:"中央値を出力する"`
	Percentile          int                   `short:"p" long:"percentile" description:"パーセンタイル値を出力する(1~100)"`
	SortedFlag          bool                  `short:"s" long:"sorted" description:"入力元データがソート済みフラグ"`
	HeaderFlag          bool                  `short:"H" long:"header" description:"ヘッダを出力する"`
	InputDelimiter      string                `short:"d" long:"indelimiter" description:"入力の区切り文字を指定" default:"\t"`
	OutputDelimiter     string                `short:"D" long:"outdelimiter" description:"出力の区切り文字を指定" default:"\t"`
	OutFile             string                `short:"o" long:"outfile" description:"出力ファイルパス"`
	SeparatableFilePath []SeparatableFilePath `short:"f" long:"fieldfilepath" description:"複数フィールド持つファイルと、その区切り位置指定(N:filepath)"`
	IgnoreHeaderRows    int                   `short:"I" long:"ignoreheader" description:"入力データヘッダを指定行無視する"`
}

type SeparatableFilePath struct {
	FieldIndex int
	FilePath   string
}

func (s *SeparatableFilePath) UnmarshalFlag(v string) error {
	const sep = ":"

	// 空文字はNG
	if strings.TrimSpace(v) == "" {
		return errors.New("not allowed empty value.")
	}

	// 区切り文字がなければ、ファイルパスとしてそのまま返す
	// インデックスは1
	if !strings.Contains(v, sep) {
		s.FieldIndex = 1
		s.FilePath = v
		return nil
	}

	// 空白を切り詰めて数値変換
	parts := strings.Split(v, sep)

	stri := strings.TrimSpace(parts[0])
	fn := parts[1]

	if stri == "" || fn == "" {
		msg := fmt.Sprintf("value is empty. index=%s filename=%s", stri, fn)
		return errors.New(msg)
	}

	i, err := strconv.Atoi(stri)
	if err != nil {
		return errors.New("expected that first values is integer that separated by a : .")
	}

	// 1未満はNG
	if i < 1 {
		msg := fmt.Sprintf("integer is over 1. input=%v", i)
		return errors.New(msg)
	}

	s.FieldIndex = i
	s.FilePath = fn

	return nil
}

func (s SeparatableFilePath) MarshalFlag() (string, error) {
	return fmt.Sprintf("%d:%s", s.FieldIndex, s.FilePath), nil
}

// OutValues はFormat関数で使用する値構造体です。
type OutValues struct {
	FileName   string
	Count      int
	Min        float64
	Max        float64
	Sum        float64
	Average    float64
	Median     float64
	Percentile float64
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

	// -f フラグがあるときはファイルパスを上書きする
	l := len(opts.SeparatableFilePath)
	if 1 <= l {
		fns := make([]string, l)
		for i, v := range opts.SeparatableFilePath {
			fns[i] = v.FilePath
		}
		args = fns
	}

	return opts, args
}

// Setup はオプションのデフォルト値をセットします。
// Count, Min, Max, Sumのいずれもfalseの場合は、すべてtrueにする。
func (o *Options) Setup() {
	if !o.CountFlag &&
		!o.MinFlag &&
		!o.MaxFlag &&
		!o.SumFlag &&
		!o.AverageFlag &&
		!o.MedianFlag &&
		o.Percentile <= 0 {
		o.CountFlag = true
		o.MinFlag = true
		o.SumFlag = true
		o.MaxFlag = true
		o.AverageFlag = true
		o.MedianFlag = true
		o.Percentile = 95
	}
	if 100 < o.Percentile {
		msg := fmt.Sprintf("warn: percentile is from 1 to 100. percentile=%d", o.Percentile)
		fmt.Fprintln(os.Stderr, msg)
		o.Percentile = 100
	}
}

// Format は出力用のデータをオプションに応じて出力ように整形する。
func Format(vs []OutValues, opts Options) []string {
	percentileHeader := fmt.Sprintf("%d%s", opts.Percentile, HeaderPercentile)

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

		if v.FileName != "" {
			setFunc(!opts.NoFileNameFlag, FileName, v.FileName)
		}
		setFunc(opts.CountFlag, HeaderCount, v.Count)
		setFunc(opts.MinFlag, HeaderMin, v.Min)
		setFunc(opts.MaxFlag, HeaderMax, v.Max)
		setFunc(opts.SumFlag, HeaderSum, v.Sum)
		setFunc(opts.AverageFlag, HeaderAverage, v.Average)
		setFunc(opts.MedianFlag, HeaderMedian, v.Median)

		setFunc(0 < opts.Percentile, percentileHeader, v.Percentile)

		maps[i] = m
	}

	lines := make([]string, 0)

	// ヘッダの連結
	// オプションがあるとセットしない
	headers := make([]string, 0)
	m := maps[0]
	for _, k := range []string{
		FileName,
		HeaderCount,
		HeaderMin,
		HeaderMax,
		HeaderSum,
		HeaderAverage,
		HeaderMedian,
		percentileHeader,
	} {
		if m[k] != "" {
			headers = append(headers, k)
		}
	}

	s := strings.Join(headers, opts.OutputDelimiter)
	if opts.HeaderFlag {
		lines = append(lines, s)
	}

	// 値の追加
	for _, m := range maps {
		values := make([]string, 0)
		for _, k := range headers {
			values = append(values, m[k])
		}
		s := strings.Join(values, opts.OutputDelimiter)
		lines = append(lines, s)
	}

	return lines
}
