package downloder

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	rangeHeaderKey        = "Range"
	acceptRangesHeaderKey = "Accept-Ranges"
)

const acceptRangeBytes = "bytes"

type fileForSort struct {
	index    int
	filePath string
}

type indexChunk struct {
	start, end uint64
}

// Downloader that implement download methods.
type Downloader struct{}

// NewDownloader generate instance of Downloader.
func NewDownloader() *Downloader {
	return &Downloader{
	}
}

// DownloadFile download file from designated URL by HTTP Request.
// In case HTTP Server accepts download using range of bytes, it do concurrency download.
func (d *Downloader) DownloadFile(
	url url.URL,
	dst string,
) error {
	res, err := http.Get(url.String())
	if err != nil {
		return d.wrappedError(err)
	}
	defer res.Body.Close()

	acceptRange := res.Header.Get(acceptRangesHeaderKey)
	switch acceptRange {
	// Concurrency download thread using range of bytes.
	case acceptRangeBytes:
		procs := runtime.NumCPU()
		err = d.downloadConcurrency(url, dst, res, procs)

	// Single download.
	default:
		err = d.downloadSingle(dst, res)
	}
	if err != nil {
		return d.wrappedError(err)
	}

	return nil
}

// Single download.
func (d *Downloader) downloadSingle(
	dst string,
	res *http.Response,
) error {
	out, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return d.wrappedError(err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return d.wrappedError(err)
	}

	return nil
}

// Concurrency download thread using range of bytes.
func (d *Downloader) downloadConcurrency(
	url url.URL,
	dst string,
	res *http.Response,
	procs int,
) error {
	cLen := res.Header.Get("Content-Length")
	size, _ := strconv.ParseUint(cLen, 10, 64)

	// chunk file size in order for chunk file download.
	chunked := d.chunkFileSize(size, uint64(procs))

	var filePaths = make([]fileForSort, 0, len(chunked))
	eg, ctx := errgroup.WithContext(context.Background())

	for index, chunk := range chunked {
		index := index + 1
		chunk := chunk

		eg.Go(func() error {
			select {
			case <-ctx.Done():
				fmt.Println("Client Closed Request.")
				return nil

			default:
				fmt.Printf("process %d-%d\n", chunk.start, chunk.end)
				filepath, err := d.downloadChunk(url, dst, chunk, index)
				if err != nil {
					return d.wrappedError(err)
				}

				filePaths = append(filePaths, fileForSort{
					index:    index,
					filePath: *filepath,
				})

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return d.wrappedError(err)
	}

	sort.SliceStable(filePaths, func(i, j int) bool {
		return filePaths[i].index < filePaths[j].index
	})

	if err := d.mergeChunkedFile(dst, filePaths); err != nil {
		return d.wrappedError(err)
	}

	return nil
}

// Download chunk by range of bytes and return file name of downloaded.
func (d *Downloader) downloadChunk(
	url url.URL,
	dst string,
	chunk indexChunk,
	index int,
) (*string, error) {
	filePath := fmt.Sprintf("%s%s%s", dst, "_", strconv.Itoa(index))
	out, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, d.wrappedError(err)
	}
	defer out.Close()

	req, _ := http.NewRequest("GET", url.String(), nil)

	byteRange := fmt.Sprintf("bytes=%d-%d", chunk.start, chunk.end)
	req.Header.Set(rangeHeaderKey, byteRange)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, d.wrappedError(err)
	}
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return nil, d.wrappedError(err)
	}

	return &filePath, nil
}

// Chunk file size and return indexed chunk.
func (d *Downloader) chunkFileSize(
	size uint64,
	chunkSize uint64,
) []indexChunk {
	splitSize := size / chunkSize

	var chunked = make([]indexChunk, 0)
	for i := uint64(0); i < size; i += splitSize {
		idx := indexChunk{
			start: i, end: i + splitSize,
		}

		if size < idx.end {
			idx.end = size
		}

		chunked = append(chunked, idx)
	}

	return chunked
}

// merge chunked file.
func (d *Downloader) mergeChunkedFile(
	dst string,
	filePaths []fileForSort,
) error {
	out, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return d.wrappedError(err)
	}
	defer out.Close()

	for _, f := range filePaths {
		func() error {
			file, err := os.Open(f.filePath)
			if err != nil {
				return d.wrappedError(err)
			}
			defer file.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				return d.wrappedError(err)
			}
			
			if err := os.Remove(file.Name()); err != nil {
				return d.wrappedError(err)
			}

			return nil
		}()
	}

	return nil
}

func (d *Downloader) wrappedError(e error) error {
	return errors.Wrap(e, e.Error())
}