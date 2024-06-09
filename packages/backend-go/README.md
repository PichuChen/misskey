目前這個部分仍然在實驗階段，由於 Pichu 的母語是華語，所以這邊先以華語進行筆記。

這部分在短中期都不會有合併回 Misskey 的計畫，使用這部分的程式碼請自行評估風險。

2024/06/04:

目前實作一半的 /api/user/show Endpoint，簡單進行 POC 的測試數據如下：
1000 次呼叫中

* NodeJS 版本:
	* 如果使用者存在的情況：6.5 秒
	* 如果使用者不存在的情況：2.8 秒
* Go 版本:
	* 如果使用者存在的情況：0.95 秒
	* 如果使用者不存在的情況：0.922 秒

測試使用的平台是跑在 i7-12700K 的 Win11 的 WSL 虛擬機器中，同時還使用了 DevContainer

* node 版本為 v20.12.2
* golang 版本為 1.22.3
* Docker 版本為 v25.0.3

測試使用的腳本在 benchmarkbackend.go 當中

目前這個部分實作大約花了 3.5 個小時進行，如果要完整完成可能還需要兩到三小時左右。
目前計算進行加速實作是否符合經濟效益的算式大致如下：

```
p # Engineer cost per hour
t # Time spent on Golang implementation
n_i # Percentage of case i in the total number of server request
c_i # Time spent on case i in NodeJS implementation
g_i # Time spent on case i in Golang implementation
n_0 = 1 - n_1 - n_2 - n_3 - ... # Percentage of not modified case in the total number of server request

X # The server cost per month
M # How many months the server will be used

p * t < X * M * (1 - (n_1 * (g_1 / c_1) + n_2 * (g_2 / c_2) + n_3 * (g_3 / c_3) + ... + n_0) )

```
舉例來說：
```
eg:
p = 100
t = 3.5
X = 100
M = 12
n_1 = 0.05
c_1 = 6.5
g_1 = 0.95
n_2 = 0.05
c_2 = 2.8
g_2 = 0.922
n_0 = 0.9

100 * 3.5 ? 100 * 12 * (1 - (0.05 * (0.95 / 6.5) + 0.05 * (0.922 / 2.8) + 0.9) )
350 ? 100 * 12 * (1 - (0.05 * 0.146 + 0.05 * 0.329 + 0.9) )
350 ? 100 * 12 * (1 - (0.073 + 0.01645 + 0.9) )
350 ? 100 * 12 * (1 - 0.98945)
350 ? 100 * 12 * 0.01055
350 ? 1.055 * 12
350 ? 12.66
```

所以當總社群規模消耗的費用在每個月 100 美元左右時是沒有經濟效益的，但是如果到 2,700 美元的時候就會開始有經濟效益。
* 這邊的 `n_i` 並沒有參考實際的狀況，實際上約有 372 個 Endpoint, 因此平均而言 `n_i` 應該是 0.00268817204，但是實際場景中該部會是完全平均分布的。

