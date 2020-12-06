package downloder

import "net/url"

type IDownloader interface {
	Download(url url.URL) error
}
