package rest

import (
	"fmt"
	"github.com/slspeek/gotube/common"
)

func videoPostProcess(in common.Video, public bool) (out common.CVideo) {
	out = common.CVideo{Id: in.Id.Hex(),
		Owner:  in.Owner,
		Name:   in.Name,
		Desc:   in.Desc,
    Public: in.Public,
		Thumbs: make([]string, len(in.Thumbs))}
	var prefix string
	if public {
		prefix = "/public/content/videos"
	} else {
		prefix = "/content/videos"
	}
	if in.BlobId != "" {
		out.Stream = fmt.Sprintf("%s/%s", prefix, out.Id)
		out.Download = fmt.Sprintf("%s/%s/download", prefix, out.Id)
	}
	for i, _ := range in.Thumbs {
		out.Thumbs[i] = fmt.Sprintf("%s/%s/thumbs/%d", prefix, out.Id, i)
	}
	return
}
