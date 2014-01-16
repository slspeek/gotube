package thumb

import (
	"github.com/slspeek/goblob"
)

type ThumbResult struct {
	blobId string
	err    error
}

func CreateThumb(bs *goblob.BlobService, inputf string, time string, size int) (target *chan ThumbResult) {
  target = new(chan ThumbResult)
  go func() {

  }();

	return
}
