package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmcvetta/napping"
	"github.com/slspeek/flowgo"
	"github.com/slspeek/gotube/common"
	"io"
	"net/http"
	"os"
)

type Config struct {
	EndPoint, Username, Password string
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Println("I need a Config file")
		return
	}
  mesg, err := put(os.Args, cfg)
  if err != nil {
    fmt.Println("Error encountered: ", err)
    return
  }
  fmt.Println("Success: ", mesg)
}

func parseArgs(args []string) (fn, name, desc string, err error) {
	l := len(args)
	if l < 2 || l > 4 {
		err = errors.New(fmt.Sprintf("usage: %s filename [name] [description]",
			args[0]))
		return
	}
	fn = args[1]
	if l > 2 {
		name = args[2]
	} else {
		name = fn
	}
	if l > 3 {
		desc = args[3]
	} else {
		desc = ""
	}
	return
}

func parseConfig() (cfg *Config, err error) {
	f, err := os.Open(".gotube.json")
	if err != nil {
		return
	}
	return parseConfigImpl(f)
}

func parseConfigImpl(r io.Reader) (cfg *Config, err error) {
	cfg = new(Config)
	dec := json.NewDecoder(r)
	if err := dec.Decode(cfg); err == io.EOF {
		err = nil
	}
	return
}

func put(args []string, cfg *Config) (mesg string, err error) {
	fn, name, desc, err := parseArgs(args)
	if err != nil {
		return
	}
	videoMap := map[string]string{ "Owner": cfg.Username, "Name": name, "Desc": desc }
	result := new(common.Video)
	var errMsg = make(map[string]interface{})
	url := cfg.EndPoint + "/api/videos"
	req, err := http.NewRequest("", "", nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(cfg.Username, cfg.Password)
  basicAuthCode := req.Header.Get("Authorization")

	r := napping.Request{
		Method:  "POST",
		Url:     url,
		Payload: &videoMap,
		Result:  &result,
		Error:   &errMsg,
		Header:  &req.Header,
	}

	r.Header.Add("Content-Type", "application/json")
	resp, err := napping.Send(&r)
	status := resp.HttpResponse().StatusCode
	if err != nil {
		return
	}
	if status != 201 {
		err = errors.New(fmt.Sprintf("Could not create Video object online: %s",
			resp.RawText()))
		return
	} 
  uploadClient := flow.NewClient(fmt.Sprintf("%s/%s/upload", url , result.Id.Hex()))
  uploadClient.Opts.Headers["Authorization"] = basicAuthCode
  err = uploadClient.UploadFile(fn)
  if err != nil {
    return
  }
  mesg = fmt.Sprintf("Uploaded %s", fn) 
	return
}
