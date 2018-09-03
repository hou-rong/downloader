package common

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"sync"
	"time"
	"log"
	"path"
	"strings"
)

var (
	mutex = &sync.Mutex{}
)

type Block struct {
	Start  int64 `json:"start"`
	Offset int64 `json:"offset"`
	End    int64 `json:"end"`
}

type DownloadInfo struct {
	StartTime  time.Time `json:"start_time"`
	Url        string    `json:"url"`
	FilePath   string    `json:"file_path"`
	InfoPath   string    `json:"info_path"`
	Total      int64     `json:"total"`
	Blocks     []Block   `json:"blocks"`
	RetryTimes int64     `json:"retry_times"`
	ErrorCount int64     `json:"error_count"`
}

func GetFileSize(url string) int64 {
	resp, err := http.Head(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.ContentLength
}

func GetFileName(url string) string {
	strList := strings.Split(url, "/")
	return strList[len(strList)-1]
}

func LoadDownloadInfo(infoPath string) DownloadInfo {
	mutex.Lock()
	f, err := os.Open(infoPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	mutex.Unlock()

	var downloadInfo DownloadInfo
	err = json.Unmarshal(byteValue, &downloadInfo)
	if err != nil {
		panic(err)
	}
	return downloadInfo
}

func SaveDownloadInfo(downloadInfo *DownloadInfo) {
	byteValue, err := json.Marshal(downloadInfo)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(downloadInfo.InfoPath); err != nil {
		os.OpenFile(downloadInfo.InfoPath, os.O_RDONLY|os.O_CREATE, 0666)
	}

	mutex.Lock()
	f, err := os.OpenFile(downloadInfo.InfoPath, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(byteValue)
	if err != nil {
		panic(err)
	}
	mutex.Unlock()
}

func DownloadLog(downloadInfo *DownloadInfo, id int64, err error) {
	downloaded := float64(0)
	for _, blocks := range downloadInfo.Blocks {
		downloaded += float64(blocks.Offset - blocks.Start)
	}
	if err != nil {
		downloadInfo.ErrorCount += 1
		log.Printf(
			"[Failed][Percent: %0.2f%%][Takes: %s][TotalRetrys: %d][Errors: %d][ErrorInfo: %s][ %v ]\n",
			100*downloaded/float64(downloadInfo.Total),
			time.Now().Sub(downloadInfo.StartTime),
			downloadInfo.RetryTimes,
			downloadInfo.ErrorCount,
			err,
			downloadInfo.Blocks[id],
		)
	} else {
		log.Printf(
			"[Finished][Percent: %0.2f%%][Takes: %s][TotalRetrys: %d][Errors: %d][ %v ]\n",
			100*downloaded/float64(downloadInfo.Total),
			time.Now().Sub(downloadInfo.StartTime),
			downloadInfo.RetryTimes,
			downloadInfo.ErrorCount,
			downloadInfo.Blocks[id],
		)
	}

	SaveDownloadInfo(downloadInfo)
}

func DownloadBlock(downloadInfo *DownloadInfo, id int64, f *os.File) (error) {
	request, err := http.NewRequest("GET", downloadInfo.Url, nil)
	if err != nil {
		panic(err)
	}

	start := downloadInfo.Blocks[id].Start
	end := downloadInfo.Blocks[id].End
	request.Header.Set(
		"Range",
		fmt.Sprintf("bytes=%d-%d", start, end),
	)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		DownloadLog(downloadInfo, id, err)
		return err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		DownloadLog(downloadInfo, id, err)
		return err
	}
	downloadInfo.Blocks[id].Offset += int64(len(buf))
	f.WriteAt(buf, start)
	DownloadLog(downloadInfo, id, nil)
	return nil
}

func Download(url string, filePath string, ext string) {
	if filePath == "" {
		filePath = path.Join(DownloadFolder, GetFileName(url))
	}
	if ext == "" {
		ext = ".downloading"
	}

	if _, err := os.Stat(filePath); err == nil {
		log.Fatalf("[File Already Exists][Path: %s]\n", filePath)
	}

	workPath := fmt.Sprintf("%s%s", filePath, ext)
	infoPath := fmt.Sprintf("%s.info.json", filePath)

	var downloadInfo DownloadInfo

	if _, err := os.Stat(infoPath); err == nil {
		downloadInfo = LoadDownloadInfo(infoPath)
	} else {
		total := GetFileSize(url)

		blockCount := total / BlockSize
		remain := total % BlockSize

		var blocks []Block
		for i := int64(0); i < blockCount; i++ {
			var start, end int64;
			start = i * BlockSize
			if i == (blockCount - 1) {
				end = (i+1)*BlockSize + remain
			} else {
				end = (i+1)*BlockSize - 1
			}
			blocks = append(blocks, Block{start, start, end})
		}

		downloadInfo = DownloadInfo{
			time.Now(),
			url,
			workPath,
			infoPath,
			total,
			blocks,
			0,
			0,
		}

		SaveDownloadInfo(&downloadInfo)

		os.OpenFile(workPath, os.O_RDONLY|os.O_CREATE, 0666)

	}

	log.Printf("[Start Download][Total: %d][Path: %s]\n", downloadInfo.Total, filePath)

	workFile, err := os.OpenFile(workPath, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var finished = false
	for !finished {
		finished = true
		for id := range downloadInfo.Blocks {
			if downloadInfo.Blocks[id].Offset-downloadInfo.Blocks[id].End >= 0 {
				continue
			}
			wg.Add(1)
			go func(downloadInfo *DownloadInfo, id int64, f *os.File) {
				defer wg.Done()
				err := DownloadBlock(downloadInfo, id, f)
				if err != nil {
					finished = false
				}
			}(&downloadInfo, int64(id), workFile)
		}
		wg.Wait()
		downloadInfo.RetryTimes += 1
	}

	log.Printf("[Download Finished][Total: %d][Takes: %s]\n", downloadInfo.Total, time.Now().Sub(downloadInfo.StartTime), )

	// 更名
	if _, err := os.Stat(filePath); err != nil {
		os.Remove(filePath)
		os.Rename(workPath, filePath)
	}

	// 删除信息文件
	if !Debug {
		if _, err := os.Stat(infoPath); err == nil {
			os.Remove(infoPath)
		}
	}
}
