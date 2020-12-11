package iprogresser

import "context"

// IProgresser is interface of Progresser.
type IProgresser interface {
	WriteProgressBar(
		ctx context.Context,
		fileSize int64,
		index int,
		isAcceptRange bool,
	)
}
