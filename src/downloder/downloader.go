package downloder

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const acceptRangesHeaderKey = "Accept-Ranges"

const acceptRangeBytes = "bytes"

type indexChunk struct {
	start, end uint64
}

type Downloader struct{
	Procs int
}

func NewDownloader() *Downloader {
	return &Downloader{
		Procs: runtime.NumCPU(),
	}
}

func (d *Downloader) HTTPDownloadFile(
	url url.URL,
	dst string,
) error {
	res, err := http.Get(url.String())
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	defer res.Body.Close()

	acceptRange := res.Header.Get(acceptRangesHeaderKey)

	switch acceptRange {
	case acceptRangeBytes:
		err = d.downloadConcurrency(url, dst, res)

	default:
		err = d.downloadSingle(dst, res)
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
		return errors.Wrap(err, err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	return nil
}

func (d *Downloader) downloadConcurrency(
	url url.URL,
	dst string,
	res *http.Response,
) error {
	out, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	defer out.Close()

	size := d.getFileSize(res)
	chunked := d.chunkFileSize(size, uint64(d.Procs))

	eg, ctx := errgroup.WithContext(context.Background())
	for _, chunk := range chunked {
		chunk := chunk
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				fmt.Println("Client Closed Request.")
				return nil

			default:
				fmt.Printf("process %d-%d\n", chunk.start, chunk.end)
				err := d.downloadChunk(url, out, chunk)
				if err != nil {
					return err
				}

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (d *Downloader) downloadChunk(
	url url.URL,
	out *os.File,
	chunk indexChunk,
) error {
	req, _ := http.NewRequest("GET", url.String(), nil)

	byteRange := fmt.Sprintf("bytes=%d-%d", chunk.start, chunk.end)
	req.Header.Set("Range", byteRange)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return errors.Wrap(err, err.Error())
	}

	return nil
}

func (d *Downloader) getFileSize(
	res *http.Response,
) uint64 {
	cLen := res.Header.Get("Content-Length")
	if cLen == "" {
		return 0
	}

	size, _ := strconv.ParseUint(cLen, 10, 64)
	return size
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