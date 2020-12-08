package idownloader

import "net/url"

type IDownloader interface {
	HTTPDownloadFile(url url.URL) error
}
