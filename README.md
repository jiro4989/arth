# arth (arithmetic)

[![Build Status](https://travis-ci.org/jiro4989/arth.svg?branch=master)](https://travis-ci.org/jiro4989/arth)

算術コマンドラインツール。

## 目的

負荷試験の結果をまとめたりするときに
最小値、最大値、平均値、中央値とかをまとめることが多い。

毎回awkで計算するのもアホらしいのでCLIにしてみた。

## できること

集計したデータの下記のデータを出力できる。

1. 件数
1. 最小値
1. 最大値
1. 合計値
1. 平均値
1. 中央値
1. パーセンタイル値

また、集計データを複数同時に並列集計することが可能。  
実行方法は「使い方/複数ファイル指定」を参照。

## インストール方法

`go get github.com/jiro4989/arth`

## 使い方

```
# num.txt
1
4
2
5
3
```

```bash
$ arth num.txt
testdata/normal_num.txt	5	1	5	15	3	3	4

# ヘッダ有り
$ arth num.txt -H
filename	count	min	max	sum	avg	median	95percentile
testdata/normal_num.txt	5	1	5	15	3	3	4
```

### 複数ファイル指定

```bash
$ arth testdata/bench.txt testdata/normal_num.txt testdata/bigdata.txt 
testdata/bench.txt	6000000	1	6000000	18000003000000	3000000.5	3000000	5700000
testdata/normal_num.txt	5	1	5	15	3	3	4
testdata/bigdata.txt	100	1	100	5050	50.5	50	95
```

### フィールド指定

`\d:filepath`と指定することで、カラム指定でファイルを読み込める。

```bash
$ arth -d , -D , -I 1 -f 2:testdata/sample.csv -f 3:testdata/sample.csv     
testdata/sample.csv,5,70,90,400,80,77,88
testdata/sample.csv,5,80,80,400,80,80,80
```

## ヘルプ

`arth -h`

    Usage:
      arth [OPTIONS]

    Application Options:
      -v, --version        バージョン情報
      -N, --nofilename     入力元ファイル名を出力しない
      -c, --count          データ数を出力する
      -n, --min            最小値を出力する
      -x, --max            最大値を出力する
      -u, --sum            合計を出力する
      -a, --avg            平均値を出力する
      -m, --median         中央値を出力する
      -p, --percentile=    パーセンタイル値を出力する(1~100)
      -s, --sorted         入力元データがソート済みフラグ
      -H, --header         ヘッダを出力する
      -d, --indelimiter=   入力の区切り文字を指定 (default: "\t")
      -D, --outdelimiter=  出力の区切り文字を指定 (default: "\t")
      -o, --outfile=       出力ファイルパス
      -f, --fieldfilepath= 複数フィールド持つファイルと、その区切り位置指定(N:filep-

                           ath)
      -I, --ignoreheader=  入力データヘッダを指定行無視する

    Help Options:
      -h, --help           Show this help message

## 仕様

### 不正なデータ

読み込んだデータに数値以外のものが混じっていた場合は集計対象から無視して
計算を続行する。その場合、データ総数(count)にも含めない。

### オプション引数

count,min,max,sum,avg,median,percentileはデフォルトですべて出力する。

ただし、上記のいずれかを指定した場合、ファイルパスとそのオプションの値のみ出力さ
れる。

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
median.sh

real	0m4.036s
user	0m10.165s
sys	0m0.099s

arth

real	0m2.292s
user	0m2.312s
sys	0m0.033s
================================
loop arth vs arth goroutine
loop arth

real	0m6.896s
user	0m6.948s
sys	0m0.114s

arth goroutine

real	0m3.399s
user	0m9.441s
sys	0m0.073s
```
