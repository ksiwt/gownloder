package main

import (
	"fmt"
	"net/url"
	"os"

	"gownloder/src/downloder"
	"gownloder/src/downloder/idownloader"
)

const URL = "http://i.imgur.com/z4d4kWk.jpg"

func main() {
	// dummy
	url, _ := url.ParseRequestURI(URL)
	dst := "./out" + url.Path

	var iDownloader idownloader.IDownloader
	iDownloader = downloder.NewDownloader()

	if err := iDownloader.DownloadFile(*url, dst); err != nil {
		fmt.Printf("Download Failed: %+v", err)

		os.Exit(1)
	}

	fmt.Println("Download Successfully!")
}
