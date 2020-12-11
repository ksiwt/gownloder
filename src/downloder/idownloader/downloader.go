package idownloader

import (
	"net/url"
)

// IDownloader is interface of Downloader.
type IDownloader interface {
	// DownloadFile download file from designated URL by HTTP Request.
	// In case HTTP Server accepts download using range of bytes, it do concurrency download.
	DownloadFile(url url.URL, dst string)  error
}
