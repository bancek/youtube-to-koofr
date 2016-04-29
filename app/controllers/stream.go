package controllers

import (
	"fmt"
	"github.com/bancek/youtube-to-koofr/app/models"
	"github.com/revel/revel"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type Stream struct {
	*revel.Controller
}

type StreamConvertResult struct {
	url string
}

func (r *StreamConvertResult) Apply(req *revel.Request, resp *revel.Response) {
	resp.Out.Header().Add("Access-Control-Allow-Origin", "*")

	resp.WriteHeader(http.StatusOK, "audio/mp3")

	logger := func(line string) {
		if revel.DevMode {
			fmt.Println(line)
		}
	}

	tmpDir, err := ioutil.TempDir("", "youtube-to-koofr")
	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	fileName, err := models.YoutubeDl(r.url, tmpDir, logger)
	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	filePath := path.Join(tmpDir, fileName)

	f, err := os.Open(filePath)

	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	_, err = io.Copy(resp.Out, f)

	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	return
}

func (c Stream) Convert(url string) revel.Result {
	return &StreamConvertResult{
		url: url,
	}
}
