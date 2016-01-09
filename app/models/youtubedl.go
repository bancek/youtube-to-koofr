package models

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func YoutubeDl(url string, destDir string, logger func(string)) (fileName string, err error) {
	cmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", "--audio-quality", "0", url)
	cmd.Dir = destDir

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		logger(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	files, err := ioutil.ReadDir(destDir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("Audio file not found")
	}

	fileName = files[0].Name()

	return fileName, nil
}
