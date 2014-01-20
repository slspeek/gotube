package thumb

import (
	"fmt"
	"github.com/slspeek/goblob"
	"github.com/slspeek/goffthumb"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type ThumbResult struct {
	BlobId string
	Index  int
	Err    error
}

func timePercentage(count int, index int) (label string) {
	base := (100 / (2 * count))
	value := base + index*(100/count)
	rounded := int(value)
	label = fmt.Sprintf("%d%%", rounded)
	log.Println("return: ", label)
	return
}

func CreateThumbs(bs *goblob.BlobService, inputf string, count int, size int) (target chan ThumbResult, err error) {
	target = make(chan ThumbResult, count)
	limit := make(chan int, 2)
	var wg sync.WaitGroup
	inputdfn, err := tempFileName()
	if err != nil {
		return
	}
	err = bs.WriteOutFile(inputf, inputdfn)
	if err != nil {
		return
	}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			limit <- 1
			blobId, err := createThumbImpl(bs, inputdfn, timePercentage(count, i), size)
			target <- ThumbResult{BlobId: blobId, Index: i, Err: err}
			<-limit
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
    os.Remove(inputdfn)
	}()
	return
}

func createThumbImpl(bs *goblob.BlobService, inputdfn string, time string, size int) (blobId string, err error) {
	thumbfn, err := tempFileName()
	if err != nil {
		return
	}
	thumbfn += ".png"
	err = ffmpeg.CreateThumb(inputdfn, thumbfn, time, size)
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
