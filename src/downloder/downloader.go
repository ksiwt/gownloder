package downloder

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const acceptRangesHeaderKey = "Accept-Ranges"

const acceptRangeBytes = "bytes"

type indexChunk struct {
	start uint64
	end   uint64
}

type Downloader struct{}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) HTTPDownloadFile(url url.URL) error {
	res, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("HTTP status error: %s", res.Status)
	}

	acceptRange := res.Header.Get(acceptRangesHeaderKey)

	// TODO: コマンドで指定する。
	path := "./out" + url.Path

	switch acceptRange {
	case acceptRangeBytes:
		// chunkSize は並列処理の同時実行数(利用可能なCPUのコア数)を基準とする。
		cpus := runtime.NumCPU()
		err = d.downloadChunk(path, res, cpus)

	default:
		err = d.downloadSingle(path, res)
	}
	if err != nil {
		return err
	}

	return nil
}

func (d *Downloader) downloadSingle(
	dst string,
	res *http.Response,
) error {
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

func (d *Downloader) downloadChunk(
	dst string,
	res *http.Response,
	chunkSize int,
) error {
	size, err := d.getFileSizeFromContentLength(res)
	if err != nil {
		return err
	}

	var (
		chunked = d.chunkFileSize(size, uint64(chunkSize))
		ch      = make(chan struct{}, chunkSize)
		eg      = errgroup.Group{}
	)

	// TODO:

	return nil
}

func (d *Downloader) chunkFileSize(
	size uint64,
	chunkSize uint64,
) []indexChunk {
	splitSize := size / chunkSize

	var chunked = make([]indexChunk, 0)
	for i := uint64(0); i < size; i += splitSize {
		idx := indexChunk{
			start: i,
			end:   i + splitSize,
		}

		if size < idx.end {
			idx.end = size
		}

		chunked = append(chunked, idx)
	}

	return chunked
}

func (d *Downloader) getFileSizeFromContentLength(
	res *http.Response,
) (uint64, error) {
	cLen := res.Header.Get("Content-Length")
	if cLen == "" {
		return 0, nil
	}

	size, err := strconv.ParseUint(cLen, 10, 64)
	if err != nil {
		return 0, errors.New("Cannot read Content-Length")
	}

	return size, nil
}
