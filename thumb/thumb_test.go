package thumb

import (
	"github.com/slspeek/goblob"
	"labix.org/v2/mgo"
	"testing"
)

func TestTimeFunction(t *testing.T) {
	count := 5
	if timePercentage(count, 0) != "10%" {
		t.Fail()
	}

	if timePercentage(count, 4) != "90%" {
		t.Fail()
	}
}

func blobService() (bs *goblob.BlobService) {
	s, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	bs = goblob.NewBlobService(s, "test", "flowfs")
	return
}

func insertBBBwebm() (id string, err error) {
  bs := blobService()
  id, err = bs.ReadFile("../test-data/big_buck_bunny.webm")
  return
}

func TestCreateThumbs(t *testing.T) {
  count := 10
  bs := blobService()
  id, err := insertBBBwebm()
  if err != nil {
    t.Fatal(err)
  }
  resultChan, err := CreateThumbs(bs, id, count, 256)
  if err != nil {
    t.Fatal(err)
  }
  allThumbs := make(map[int]ThumbResult)
	for i := 0; i < count; i++ {
		result := <-resultChan
		allThumbs[result.Index] = result
	}
	for i := 0; i < count; i++ {
		if err = allThumbs[i].Err; err != nil {
      t.Fail()
		}
	}
}
