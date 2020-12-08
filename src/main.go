package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"gownloder/src/downloder"
	"gownloder/src/downloder/idownloader"
)

const URL = "http://i.imgur.com/z4d4kWk.jpg"

func main() {
	// dummy
	url, _ := url.ParseRequestURI(URL)


	var iDownloader idownloader.IDownloader
	iDownloader = downloder.NewDownloader()
	if err := iDownloader.HTTPDownloadFile(*url); err != nil {
		err := errors.Wrap(err, "Download Failed")
		fmt.Println(err)

		os.Exit(1)
	}

	fmt.Println("Download Successfully!")
}
