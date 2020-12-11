package idownloader

import (
	"net/url"
)

type IDownloader interface {
	DownloadFile(url url.URL, dst string)  error
}
