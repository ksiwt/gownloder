package progresser

import (
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
	fileSize int64,
) {
	bar := pb.Start64(fileSize)
	bar.Increment()
	bar.Set(pb.Bytes, true)
	bar.SetMaxWidth(100)

	bar.Start()
	bar.Finish()
}
