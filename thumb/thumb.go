package thumb

import (
	"fmt"
	"github.com/slspeek/goblob"
	"github.com/slspeek/goffthumb"
	"io/ioutil"
	"os"
)

type ThumbResult struct {
	BlobId string
	Err    error
}

func CreateThumbs(bs *goblob.BlobService, inputf string, count int, size int) (target chan ThumbResult, tempfn string, err error) {
	target = make(chan ThumbResult, count)
	inputdfn, err := tempFileName()
  if err != nil {
    return
  }
  tempfn = inputdfn
	inputfh, err := bs.Open(inputf)
	if err != nil {
		return
	}
	err = goblob.WriteFile(inputdfn, inputfh)
	if err != nil {
		return
	}
	for time := 100 / (2 * count); time <= 100; time += 100 / count {
		out := int(time)
		go func() {
			blobId, err := createThumbImpl(bs, inputdfn, fmt.Sprintf("%v%%", out), size)
			target <- ThumbResult{BlobId: blobId, Err: err}
		}()
	}
  
	return
}

func createThumbImpl(bs *goblob.BlobService, inputdfn string, time string, size int) (blobId string, err error) {
	thumbfn, err := tempFileName()
  if err != nil {
    return
  }
  thumbfn += ".png"
	err = ffmpeg.CreateThumb(inputdfn, thumbfn , time, size)
	if err != nil {
		return
	}

	blobId, err = bs.ReadFile(thumbfn)
	if err != nil {
		return
	}
	os.Remove(thumbfn)
	return
}

func tempFileName() (path string, err error) {
	f, err := ioutil.TempFile("", "goffthumb")
	if err != nil {
		return
	}
	path = f.Name()
  err = f.Close()
  if err != nil {
    return
  }
  os.Remove(path)
	return
}
