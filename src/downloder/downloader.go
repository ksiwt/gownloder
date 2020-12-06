package downloder

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

const acceptRangesHeaderKey = "Accept-Ranges"

const (
	acceptRangeBytes = "bytes"
	acceptRangeNone = "none"
)

const defaultChunkSize = 100

type Downloader struct {}

func newDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(url url.URL) error {
	res, err := http.Get(url.String())
	if err != nil {
		return errors.Wrapf(err, "Failed HTTP GET URL:%s" , url.String())
	}
	defer res.Body.Close()

	acceptRange := res.Header.Get(acceptRangesHeaderKey)

	// TODO:
	var outDir = ""

	var downloadErr error
	switch acceptRange {
	case acceptRangeNone:
		downloadErr = d.download(url, outDir)
	case acceptRangeBytes:
		downloadErr = d.downloadChunk()
	default:
		msg := fmt.Sprintf("Invalid Accept-Ranges header: %s", acceptRange)
		return errors.New(msg)
	}
	if downloadErr!= nil {
		return errors.Wrapf(err, "Failed download %s", url.Path)
	}

	fmt.Println("Download Successfully !")
	return nil
}

func (d *Downloader) download(
	url url.URL,
	dst string,
) error {
	res, err := http.Get(url.String())
	if err != nil {
		return errors.Wrapf(err, "Failed HTTP GET URL:%s" , url.String())
	}
	defer res.Body.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	return nil
}

// TODO
func (d *Downloader) downloadChunk(
) error {
	return nil
}
