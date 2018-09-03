# Downloader

## 简介

支持并发断点续传的下载器（练手项目）

## 使用方法及结果

```go
package main

import "github.com/hou-rong/downloader/pkg/common"

func main() {

	downloadUrls := []string{
		"https://s3.amazonaws.com/expedia-static-feed/United+States+(.com)_Merchant_All.csv.gz",
		"https://www.python.org/ftp/python/3.5.6/Python-3.5.6.tar.xz",
		"https://download.calibre-ebook.com/3.30.0/calibre-3.30.0.dmg",
	}

	for _, downloadUrl := range downloadUrls {
		common.Download(
			downloadUrl,
			"",
			"",
		)
	}
}
```
**output**
    
    2018/09/03 23:06:41 [Start Download][Total: 23521935][Path: /Users/hourong/Downloads/United+States+(.com)_Merchant_All.csv.gz]
    2018/09/03 23:06:42 [Finished][Percent: 0.22%][Takes: 1.343669617s][TotalRetrys: 0][Errors: 0][ {0 51200 51199} ]
    2018/09/03 23:06:42 [Finished][Percent: 0.44%][Takes: 1.609520355s][TotalRetrys: 0][Errors: 0][ {14694400 14745600 14745599} ]
    2018/09/03 23:06:43 [Finished][Percent: 0.65%][Takes: 1.840226844s][TotalRetrys: 0][Errors: 0][ {1792000 1843200 1843199} ]
    2018/09/03 23:06:43 [Finished][Percent: 0.87%][Takes: 1.877752259s][TotalRetrys: 0][Errors: 0][ {6297600 6348800 6348799} ]
    2018/09/03 23:06:43 [Finished][Percent: 1.09%][Takes: 1.894185019s][TotalRetrys: 0][Errors: 0][ {15257600 15308800 15308799} ]
    ......
    2018/09/03 23:06:47 [Finished][Percent: 99.13%][Takes: 6.647982914s][TotalRetrys: 0][Errors: 0][ {21657600 21708800 21708799} ]
    2018/09/03 23:06:48 [Finished][Percent: 99.35%][Takes: 6.963831809s][TotalRetrys: 0][Errors: 0][ {10188800 10240000 10239999} ]
    2018/09/03 23:06:48 [Finished][Percent: 99.56%][Takes: 7.302839903s][TotalRetrys: 0][Errors: 0][ {4096000 4147200 4147199} ]
    2018/09/03 23:06:49 [Finished][Percent: 99.78%][Takes: 7.789052347s][TotalRetrys: 0][Errors: 0][ {14592000 14643200 14643199} ]
    2018/09/03 23:06:50 [Finished][Percent: 100.00%][Takes: 9.000952727s][TotalRetrys: 0][Errors: 0][ {3379200 3430400 3430399} ]
    2018/09/03 23:06:50 [Download Finished][Total: 23521935][Takes: 9.001497618s]
    
## Done & ToDo

- [x] 断点续传
- [x] 简易监控
- [ ] 参数启动
- [ ] 错误处理机制
- [ ] 独立监控功能
- [ ] 多任务
- [ ] ......