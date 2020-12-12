package iprogresser

// IProgresser is interface of Progresser.
type IProgresser interface {
	WriteProgressBar(
		fileSize int64,
	)
}
