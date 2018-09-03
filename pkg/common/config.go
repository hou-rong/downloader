package common

import (
	"os/user"
	"log"
	"path"
)

var (
	DownloadFolder       = ""
	BlockSize      int64 = 50 * 1024
	Debug                = false
)

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	DownloadFolder = path.Join(usr.HomeDir, "Downloads")
}
