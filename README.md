# arth (arithmetic)

算術コマンドラインツール。

## 目的

負荷試験の結果をまとめたりするときに
最小値、最大値、平均値、中央値とかをまとめることが多い。

毎回awkで計算するのもアホらしいのでCLIにしてみた。

## インストール方法

`go get github.com/jiro4989/arth`

## 使い方

```time.list
1
4
2
5
3
```

```bash
arth time.list
# 出力
# filename	count	min	max	sum	avg
# time.list	5	1	5	15	3
```

### 複数ファイル指定

```bash
$ arth -m testdata/bench.txt testdata/normal_num.txt testdata/bigdata.txt 
filename	count	min	max	sum	avg	median
testdata/bench.txt	6000000	1	6000000	18000003000000	3000000.5	3000000
testdata/normal_num.txt	5	1	5	15	3	3
testdata/bigdata.txt	100	1	100	5050	50.5	50
```

## ヘルプ

`arth -h`

    Usage:
      arth [OPTIONS]

    Application Options:
      -v, --version     バージョン情報
          --nofilename  入力元ファイル名を出力しない
          --count       データ数を出力する
          --min         最小値を出力する
          --max         最大値を出力する
          --sum         合計を出力する
          --avg         平均値を出力する
      -m, --median      中央値を出力する
      -s, --sorted      入力元データがソート済みフラグ
      -n, --noheader    ヘッダを出力しない
      -d, --delimiter=  出力時の区切り文字を指定 (default: "\t")
      -o, --outfile=    出力ファイルパス

    Help Options:
      -h, --help        Show this help message

## 仕様

### 不正なデータ

読み込んだデータに数値以外のものが混じっていた場合は集計対象から無視して
計算を続行する。その場合、データ総数(count)にも含めない。

### オプション引数

count,min,max,sum,avgはデフォルトで出力する。

ただし、上記のいずれもオプション引数で指定しない場合のみ上記がデフォルトで出力さ
れる。

つまり、上記の5つのうち、1つでも意図的に指定すると、他の4つが出力されなくなる。

### 通常例

```bash
$ arth testdata/bigdata.txt
filename	count	min	max	sum	avg
testdata/bigdata.txt	100	1	100	5050	50.5
```

```bash
$ arth testdata/bigdata.txt --count
filename	count
testdata/bigdata.txt	100
```

```bash
$ arth testdata/bigdata.txt --count --sum
filename	count	sum
testdata/bigdata.txt	100	5050
```

```bash
$ arth testdata/bigdata.txt -m
filename	count	min	max	sum	avg	median
testdata/bigdata.txt	100	1	100	5050	50.5	50
```

## 開発方法

```
make deps
make
```

## 処理速度

ベンチマーク用のスクリプトを実行した結果。

```bash
$ bash script/bench.sh
median.sh vs arth

real	0m4.498s
user	0m10.668s
sys	0m0.162s

real	0m2.447s
user	0m2.466s
sys	0m0.040s
================================
arth -m goroutine vs loop arth -m

real	0m3.472s
user	0m9.636s
sys	0m0.127s

real	0m14.631s
user	0m14.763s
sys	0m0.214s
```
