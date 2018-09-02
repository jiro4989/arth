package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"

	"github.com/jiro4989/arth/internal/options"
	arthio "github.com/jiro4989/arth/io"
	arthmath "github.com/jiro4989/arth/math"
)

// エラー出力ログ
var logger = log.New(os.Stderr, "", 0)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	// オプション引数の解析
	opts, args := options.Parse(Version)

	// 入力データの処理
	ovs, err := processInput(args, opts)
	if err != nil {
		panic(err)
	}

	// 出力用に整形
	lines := options.Format(ovs, opts)

	// 標準出力、あるいはファイル出力
	if err := out(lines, opts); err != nil {
		panic(err)
	}
}

// processInput は引数、オプションを判定して計算し、出力する文字列を生成する。
// 引数指定がない場合は標準入力を受け取る
// 引数指定がある場合はファイル名としてファイル読み込みを実施
func processInput(args []string, opts options.Options) ([]options.OutValues, error) {
	if len(args) < 1 {
		return processStdin(opts)
	}

	return processMultiInput(args, opts), nil
}

// processStdin は標準入力のデータを処理する。
func processStdin(opts options.Options) ([]options.OutValues, error) {
	r := os.Stdin
	ov, err := calcOutValues(r, opts)
	if err != nil {
		return nil, err
	}
	return []options.OutValues{ov}, nil
}

// indexedFileName は処理し始めた順番を保持するファイル名。
type indexedFileName struct {
	index    int
	fileName string
}

// processMultiInput は複数の入力ファイルを処理する。
// CPUの数だけワーカースレッドを起動し、並列でデータを処理する。
// FIXME goroutineの途中にエラーが発生してもエラーを返さない。
// ログ出力はするが
func processMultiInput(args []string, opts options.Options) []options.OutValues {
	var wg sync.WaitGroup
	q := make(chan indexedFileName, len(args))

	// CPUの数だけワーカースレッドを起動
	// 並列でファイルを開いて処理し、出力データ配列に追加する
	ovs := make([]options.OutValues, len(args))
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func(ovs []options.OutValues) {
			defer wg.Done()
			for {
				// 入力ファイル名を受け取る
				ifn, ok := <-q
				if !ok {
					return
				}

				fn := ifn.fileName
				ov, err := arthio.WithOpen(fn, func(r io.Reader) (options.OutValues, error) {
					return calcOutValues(r, opts)
				})
				if err != nil {
					// 処理を計測してほしいのでpanicしない
					logger.Println(err)
				}

				i := ifn.index
				ovs[i] = ov
			}
		}(ovs)
	}

	// 処理対象のファイルパスをキューに送信
	for i, fn := range args {
		ifn := indexedFileName{
			index:    i,
			fileName: fn,
		}
		q <- ifn
	}
	close(q)
	wg.Wait()

	return ovs
}

// calcOutValues は入力から出力データを計算する。
// オプションMedianFlagが存在するとき、ソートとソートデータの保持により
// メモリ消費と計算時間が増加する。
// オプションSordedFlagが存在するとき、入力がすでにソート済みとして
// ソート処理をスキップする。
func calcOutValues(r io.Reader, opts options.Options) (options.OutValues, error) {
	ov := options.OutValues{} // 出力データ
	ns := make([]float64, 0)  // 読み込んだ数値配列
	var err error
	ov.Count, ov.Min, ov.Max, ov.Sum, ov.Average, ns, err = arthmath.MinMaxSumAvg(r, opts.MedianFlag)
	if err != nil {
		return ov, err
	}

	if opts.MedianFlag {
		// SordedFlagがなければ、ソートを実行
		// SordedFlagがあれば、ソート済みとしてソートはスキップ(高速化)
		if !opts.SordedFlag {
			sort.Float64s(ns)
		}
		ov.Median = arthmath.Median(ns)
	}
	return ov, nil
}

// out は行配列をオプションに応じて出力する。
// 出力先ファイルが指定されていなければ標準出力する。
// 指定がアレばファイル出力する。
func out(lines []string, opts options.Options) error {
	if opts.OutFile == "" {
		for _, v := range lines {
			fmt.Println(v)
		}
		return nil
	}

	return arthio.WriteFile(opts.OutFile, lines)
}
