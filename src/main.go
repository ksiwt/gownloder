package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"

	"gownloder/src/downloder"
	"gownloder/src/downloder/idownloader"
	"gownloder/src/progresser"
	"gownloder/src/progresser/iprogresser"
)

func main() {
	u, d := parseFlag()
	url, _ := url.ParseRequestURI(u)

	_, file := path.Split(url.Path)
	dst := fmt.Sprintf("%s%s", d, file)

	var iProgresser iprogresser.IProgresser
	iProgresser = progresser.NewProgresser()

	var iDownloader idownloader.IDownloader
	iDownloader = downloder.NewDownloader(iProgresser)

	fmt.Println("URL:", url)
	fmt.Printf("File Download to %s \n", dst)
	if err := iDownloader.DownloadFile(*url, dst); err != nil {
		fmt.Printf("\nDownload Failed: %+v", err)

		os.Exit(1)
	}

	fmt.Println("\nDownload Successfully!")
}

func parseFlag() (url string, dst string) {
	var (
		u = flag.String("u", "", "URL of download file.")
		d = flag.String("d", "", "Destination of downloaded file.")
	)

	flag.Parse()

	if *u == "" {
		fmt.Println("Flag -u can not be null")
		os.Exit(1)
	}

	if *d == "" {
		fmt.Printf("Flag -d can not be null ")
		os.Exit(1)
	}

	return *u, *d
}