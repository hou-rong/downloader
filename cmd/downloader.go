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
