package medias

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func MergeAudioVideo(videoFile, audioFile *os.File) (err error) {
	if _, err := videoFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err := audioFile.Seek(0, 0); err != nil {
		return err
	}

	videoName := videoFile.Name()
	tempOutput := videoName + ".tmp" + filepath.Ext(videoName)
	if len(tempOutput) > 255 {
		tempOutput = tempOutput[255:]
	}
	defer func() {
		err = os.Remove(tempOutput)
		if err != nil {
			logrus.Errorf("Failed to remove temporary file: %v", err)
		}
		err = os.Remove(audioFile.Name())
		if err != nil {
			logrus.Errorf("Failed to remove audio file: %v", err)
		}
	}()

	cmd := exec.Command("ffmpeg",
		"-i", videoName,
		"-i", audioFile.Name(),
		"-c", "copy",
		"-shortest",
		"-y", tempOutput,
	)

	if err = cmd.Run(); err != nil {
		return err
	}

	tempFile, err := os.Open(tempOutput)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	if _, err = io.Copy(videoFile, tempFile); err != nil {
		return err
	}

	return err
}
