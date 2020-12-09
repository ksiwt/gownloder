package idownloader

import "net/url"

type IDownloader interface {
	HTTPDownloadFile(url url.URL, dst string) error
}
