package models

import (
	"github.com/koofr/go-koofrclient"
	"io/ioutil"
	"os"
	"path"
)

func Convert(url string, koofr *koofrclient.KoofrClient, logger func(string)) (shortUrl string, err error) {
	tmpDir, err := ioutil.TempDir("", "youtube-to-koofr")
	if err != nil {
		return "", err
	}

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	fileName, err := YoutubeDl(url, tmpDir, logger)
	if err != nil {
		return "", err
	}

	filePath := path.Join(tmpDir, fileName)

	logger("Uploading to Koofr...")

	shortUrl, err = KoofrUpload(koofr, filePath, fileName)
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}
