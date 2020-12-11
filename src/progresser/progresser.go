package progresser

import (
	"context"
	"time"

	"github.com/cheggaaa/pb/v3"
)

// Progresser implement progressbar methods.
type Progresser struct {}

// NewProgresser generate instance of Progresser.
func NewProgresser() *Progresser {
	return &Progresser{}
}

// WriteProgressBar write progress bar.
func (p *Progresser) WriteProgressBar(
	ctx context.Context,
	fileSize int64,
	index int,
	isAcceptRange bool,
) {
	// TODO:
	bar := pb.Start64(fileSize)
	bar.Increment()
	bar.SetRefreshRate(time.Second)
	bar.Set(pb.Bytes, true)
	bar.SetMaxWidth(100)

	bar.Start()
	bar.Finish()
}
