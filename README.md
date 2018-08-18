# arth (arithmetic)

算術コマンドラインツール。

## 目的

負荷試験の結果をまとめたりするときに
最小値、最大値、平均値、中央値とかをまとめることが多い。

毎回awkで計算するのもアホらしいのでCLIにしてみた。

## 使い方

### ヘルプ

`arth -h`

```
```

### 仕様例

```time.list
1
4
2
5
3
```

```bash
arim time.list
# 出力
# count	sum	min	max	avg	median
# 5	15	1	5	3	3
```

