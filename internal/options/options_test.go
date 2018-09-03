package options

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestUnmarshalFlagData struct {
	in  string
	out SeparatableFilePath
}

func TestUnmarshalFlag(t *testing.T) {
	// 正常系
	tds := []TestUnmarshalFlagData{
		TestUnmarshalFlagData{ // 期待する指定
			in: "1:foobar1.txt",
			out: SeparatableFilePath{
				FieldIndex: 1,
				FilePath:   "foobar1.txt",
			},
		},
		TestUnmarshalFlagData{ // インデックス指定なしでもよい
			in: "foobar2.txt",
			out: SeparatableFilePath{
				FieldIndex: 1,
				FilePath:   "foobar2.txt",
			},
		},
		TestUnmarshalFlagData{ // トリムしない
			in: "   foobar3.txt   ",
			out: SeparatableFilePath{
				FieldIndex: 1,
				FilePath:   "   foobar3.txt   ",
			},
		},
		TestUnmarshalFlagData{ // 空白があってもよい。ファイルパスはトリムしない
			in: " 2 : foobar4.txt ",
			out: SeparatableFilePath{
				FieldIndex: 2,
				FilePath:   " foobar4.txt ",
			},
		},
	}
	for _, v := range tds {
		s := SeparatableFilePath{}
		err := s.UnmarshalFlag(v.in)
		assert.NoError(t, err)
		assert.Equal(t, v.out, s)
	}

	// 異常系
	tds = []TestUnmarshalFlagData{
		TestUnmarshalFlagData{ // 小数
			in: "1.5:foobar.txt",
		},
		TestUnmarshalFlagData{ // 0値
			in: "0:foobar.txt",
		},
		TestUnmarshalFlagData{ // 負数
			in: "-1:foobar.txt",
		},
		TestUnmarshalFlagData{ // 数値なし
			in: ":foobar.txt",
		},
		TestUnmarshalFlagData{ // ファイルパスなし
			in: "1:",
		},
		TestUnmarshalFlagData{ // 両方なし
			in: ":",
		},
		TestUnmarshalFlagData{ // 空まじり
			in: " : ",
		},
		TestUnmarshalFlagData{ // 無
			in: "",
		},
	}
	for _, v := range tds {
		t.Log("<" + v.in + ">")
		s := SeparatableFilePath{}
		err := s.UnmarshalFlag(v.in)
		assert.Error(t, err)
	}
}

type testdata struct {
	in     Options
	expect Options
}

type TestParseData struct {
	args    []string
	outopts Options
	outargs []string
}

func TestParse(t *testing.T) {
	tds := []TestParseData{
		TestParseData{
			args: []string{
				"main.go",
			},
			outopts: Options{
				CountFlag:      true,
				MinFlag:        true,
				MaxFlag:        true,
				SumFlag:        true,
				AverageFlag:    true,
				MedianFlag:     false,
				SortedFlag:     false,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
			outargs: []string{},
		},
		TestParseData{
			args: []string{
				"main.go",
				"--count",
			},
			outopts: Options{
				CountFlag:      true,
				MinFlag:        false,
				MaxFlag:        false,
				SumFlag:        false,
				AverageFlag:    false,
				MedianFlag:     false,
				SortedFlag:     false,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
			outargs: []string{},
		},
		TestParseData{
			args: []string{
				"main.go",
				"--min",
				"testdata",
			},
			outopts: Options{
				CountFlag:      false,
				MinFlag:        true,
				MaxFlag:        false,
				SumFlag:        false,
				AverageFlag:    false,
				MedianFlag:     false,
				SortedFlag:     false,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
			outargs: []string{"testdata"},
		},
		TestParseData{ // -f で上書きされる
			args: []string{
				"main.go",
				"testdata",
				"-f",
				"sample.txt",
				"-f",
				"2:sample2.txt",
			},
			outopts: Options{
				CountFlag:      true,
				MinFlag:        true,
				MaxFlag:        true,
				SumFlag:        true,
				AverageFlag:    true,
				MedianFlag:     false,
				SortedFlag:     false,
				HeaderFlag:     false,
				InputDelimiter: "\t",
				SeparatableFilePath: []SeparatableFilePath{
					SeparatableFilePath{
						FieldIndex: 1,
						FilePath:   "sample.txt",
					},
					SeparatableFilePath{
						FieldIndex: 2,
						FilePath:   "sample2.txt",
					},
				},
			},
			outargs: []string{
				"sample.txt",
				"sample2.txt",
			},
		},
	}
	for _, v := range tds {
		os.Args = v.args
		opts, args := Parse("1.0.0")
		assert.Equal(t, v.outopts.CountFlag, opts.CountFlag)
		assert.Equal(t, v.outopts.MinFlag, opts.MinFlag)
		assert.Equal(t, v.outopts.MaxFlag, opts.MaxFlag)
		assert.Equal(t, v.outopts.SumFlag, opts.SumFlag)
		assert.Equal(t, v.outopts.AverageFlag, opts.AverageFlag)
		assert.Equal(t, v.outopts.MedianFlag, opts.MedianFlag)
		assert.Equal(t, v.outopts.SortedFlag, opts.SortedFlag)
		assert.Equal(t, v.outopts.HeaderFlag, opts.HeaderFlag)
		assert.Equal(t, v.outopts.InputDelimiter, opts.InputDelimiter)
		assert.Equal(t, v.outargs, args)
	}
}

func TestSetup(t *testing.T) {
	tds := []testdata{
		testdata{
			in: Options{
				CountFlag:      false,
				MinFlag:        false,
				MaxFlag:        false,
				SumFlag:        false,
				AverageFlag:    false,
				MedianFlag:     false,
				SortedFlag:     false,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
			expect: Options{
				CountFlag:      true,
				MinFlag:        true,
				MaxFlag:        true,
				SumFlag:        true,
				AverageFlag:    true,
				MedianFlag:     false,
				SortedFlag:     false,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
		},
		testdata{
			in: Options{
				CountFlag:      true,
				MinFlag:        false,
				MaxFlag:        false,
				SumFlag:        false,
				AverageFlag:    false,
				MedianFlag:     false,
				SortedFlag:     false,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
			expect: Options{
				CountFlag:      true,
				MinFlag:        false,
				MaxFlag:        false,
				SumFlag:        false,
				AverageFlag:    false,
				MedianFlag:     false,
				SortedFlag:     false,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
		},
		testdata{
			in: Options{
				CountFlag:      false,
				MinFlag:        true,
				MaxFlag:        false,
				SumFlag:        true,
				AverageFlag:    false,
				MedianFlag:     false,
				SortedFlag:     true,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
			expect: Options{
				CountFlag:      false,
				MinFlag:        true,
				MaxFlag:        false,
				SumFlag:        true,
				AverageFlag:    false,
				MedianFlag:     false,
				SortedFlag:     true,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
		},
		testdata{
			in: Options{
				CountFlag:      true,
				MinFlag:        true,
				MaxFlag:        true,
				SumFlag:        true,
				AverageFlag:    true,
				MedianFlag:     false,
				SortedFlag:     true,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
			expect: Options{
				CountFlag:      true,
				MinFlag:        true,
				MaxFlag:        true,
				SumFlag:        true,
				AverageFlag:    true,
				MedianFlag:     false,
				SortedFlag:     true,
				HeaderFlag:     false,
				InputDelimiter: "\t",
			},
		},
	}
	for _, v := range tds {
		v.in.Setup()
		assert.Equal(t, v.expect.CountFlag, v.in.CountFlag)
		assert.Equal(t, v.expect.MinFlag, v.in.MinFlag)
		assert.Equal(t, v.expect.MaxFlag, v.in.MaxFlag)
		assert.Equal(t, v.expect.SumFlag, v.in.SumFlag)
		assert.Equal(t, v.expect.AverageFlag, v.in.AverageFlag)
		assert.Equal(t, v.expect.MedianFlag, v.in.MedianFlag)
		assert.Equal(t, v.expect.SortedFlag, v.in.SortedFlag)
		assert.Equal(t, v.expect.HeaderFlag, v.in.HeaderFlag)
		assert.Equal(t, v.expect.InputDelimiter, v.in.InputDelimiter)
	}
}

type TestFormatData struct {
	ovs  []OutValues
	opts Options
	out  []string
}

func TestFormat(t *testing.T) {
	tds := []TestFormatData{
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:      2,
					Min:        1.0,
					Max:        2.0,
					Sum:        3.0,
					Average:    1.5,
					Median:     1.0,
					Percentile: 4.0,
				},
				OutValues{
					Count:      100,
					Min:        1.0,
					Max:        2.0,
					Sum:        3.0,
					Average:    1.5,
					Median:     1.0,
					Percentile: 95.0,
				},
			},
			opts: Options{
				CountFlag:       true,
				MinFlag:         true,
				MaxFlag:         true,
				SumFlag:         true,
				AverageFlag:     true,
				MedianFlag:      false,
				SortedFlag:      false,
				HeaderFlag:      false,
				Percentile:      95,
				InputDelimiter:  "\t",
				OutputDelimiter: "\t",
			},
			out: []string{
				"2\t1\t2\t3\t1.5\t4",
				"100\t1\t2\t3\t1.5\t95",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:      2,
					Min:        1.0,
					Max:        2.0,
					Sum:        3.0,
					Average:    1.5,
					Median:     1.0,
					Percentile: 4.0,
				},
				OutValues{
					Count:      100,
					Min:        1.0,
					Max:        2.0,
					Sum:        3.0,
					Average:    1.5,
					Median:     1.0,
					Percentile: 95.0,
				},
			},
			opts: Options{
				CountFlag:       true,
				MinFlag:         true,
				MaxFlag:         true,
				SumFlag:         true,
				AverageFlag:     true,
				MedianFlag:      true,
				SortedFlag:      false,
				HeaderFlag:      false,
				Percentile:      95,
				InputDelimiter:  "\t",
				OutputDelimiter: ",",
			},
			out: []string{
				"2,1,2,3,1.5,1,4",
				"100,1,2,3,1.5,1,95",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
				OutValues{
					Count:   100,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:       true,
				MinFlag:         true,
				MaxFlag:         true,
				SumFlag:         true,
				AverageFlag:     true,
				MedianFlag:      false,
				SortedFlag:      false,
				HeaderFlag:      false,
				InputDelimiter:  "\t",
				OutputDelimiter: "\t",
			},
			out: []string{
				"2\t1\t2\t3\t1.5",
				"100\t1\t2\t3\t1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					FileName: "foo.txt",
					Count:    2,
					Min:      1.0,
					Max:      2.0,
					Sum:      3.0,
					Average:  1.5,
					Median:   1.0,
				},
				OutValues{
					FileName: "bar.txt",
					Count:    100,
					Min:      1.0,
					Max:      2.0,
					Sum:      3.0,
					Average:  1.5,
					Median:   1.0,
				},
			},
			opts: Options{
				CountFlag:       true,
				MinFlag:         true,
				MaxFlag:         true,
				SumFlag:         true,
				AverageFlag:     true,
				MedianFlag:      false,
				SortedFlag:      false,
				HeaderFlag:      false,
				InputDelimiter:  "\t",
				OutputDelimiter: "\t",
			},
			out: []string{
				"foo.txt\t2\t1\t2\t3\t1.5",
				"bar.txt\t100\t1\t2\t3\t1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					FileName: "foo.txt",
					Count:    2,
					Min:      1.0,
					Max:      2.0,
					Sum:      3.0,
					Average:  1.5,
					Median:   1.0,
				},
				OutValues{
					FileName: "bar.txt",
					Count:    100,
					Min:      1.0,
					Max:      2.0,
					Sum:      3.0,
					Average:  1.5,
					Median:   1.0,
				},
			},
			opts: Options{
				NoFileNameFlag:  true,
				CountFlag:       true,
				MinFlag:         true,
				MaxFlag:         true,
				SumFlag:         true,
				AverageFlag:     true,
				MedianFlag:      false,
				SortedFlag:      false,
				HeaderFlag:      false,
				InputDelimiter:  "\t",
				OutputDelimiter: "\t",
			},
			out: []string{
				"2\t1\t2\t3\t1.5",
				"100\t1\t2\t3\t1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:       true,
				MinFlag:         true,
				MaxFlag:         true,
				SumFlag:         true,
				AverageFlag:     true,
				MedianFlag:      false,
				SortedFlag:      false,
				HeaderFlag:      false,
				InputDelimiter:  "\t",
				OutputDelimiter: "\t",
			},
			out: []string{
				"2\t1\t2\t3\t1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:       true,
				MinFlag:         true,
				MaxFlag:         true,
				SumFlag:         true,
				AverageFlag:     true,
				MedianFlag:      false,
				SortedFlag:      false,
				HeaderFlag:      false,
				InputDelimiter:  ",",
				OutputDelimiter: ",",
			},
			out: []string{
				"2,1,2,3,1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:       true,
				MinFlag:         true,
				MaxFlag:         true,
				SumFlag:         true,
				AverageFlag:     true,
				MedianFlag:      false,
				SortedFlag:      false,
				HeaderFlag:      true,
				InputDelimiter:  ",",
				OutputDelimiter: ",",
			},
			out: []string{
				"count,min,max,sum,avg",
				"2,1,2,3,1.5",
			},
		},
		TestFormatData{
			ovs: []OutValues{
				OutValues{
					Count:   2,
					Min:     1.0,
					Max:     2.0,
					Sum:     3.0,
					Average: 1.5,
					Median:  1.0,
				},
			},
			opts: Options{
				CountFlag:       false,
				MinFlag:         true,
				MaxFlag:         false,
				SumFlag:         true,
				AverageFlag:     true,
				MedianFlag:      false,
				SortedFlag:      false,
				HeaderFlag:      true,
				InputDelimiter:  ",",
				OutputDelimiter: ",",
			},
			out: []string{
				"min,sum,avg",
				"1,3,1.5",
			},
		},
	}

	for _, v := range tds {
		assert.Equal(t, v.out, Format(v.ovs, v.opts))
	}
}
