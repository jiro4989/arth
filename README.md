# arth (arithmetic)

算術コマンドラインツール。

## 目的

負荷試験の結果をまとめたりするときに
最小値、最大値、平均値、中央値とかをまとめることが多い。

毎回awkで計算するのもアホらしいのでCLIにしてみた。

## 使い方

### ヘルプ

`arth -h`

    Usage:
      arth [OPTIONS]

    Application Options:
          --count      データ数を出力する
          --min        最小値を出力する
          --max        最大値を出力する
          --sum        合計を出力する
          --avg        平均値を出力する
      -m, --median     中央値を出力する
      -s, --sorted     入力元データがソート済みフラグ
      -n, --noheader   ヘッダを出力しない
          --separator= 出力時のセパレータを指定 (default: "\t")
      -o, --outfile=   出力ファイルパス

    Help Options:
      -h, --help       Show this help message

### 仕様例

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
# count	min	max	sum	avg
# 5	1	5	15	3
```

## 仕様

count,min,max,sum,avgはデフォルトで出力する。

ただし、上記のいずれもオプション引数で指定しない場合のみ上記がデフォルトで出力さ
れる。

つまり、上記の5つのうち、1つでも意図的に指定すると、他の4つが出力されなくなる。

### 通常例

```bash
$ arth testdata/bigdata.txt
count	min	max	sum	avg
0	1	100	5050	50.5
```

```bash
$ arth testdata/bigdata.txt --count
count
0
```

```bash
$ arth testdata/bigdata.txt --count --sum
count	sum
0	5050
```

```bash
$ arth testdata/bigdata.txt -m
count	min	max	sum	avg	median
100	1	100	5050	50.5	50
```

