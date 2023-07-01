package opus

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"layeh.com/gopus"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

func Encode(src io.Reader, chunks chan []byte, done chan struct{}, speaking <-chan bool) error {
	run := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	run.Stdin = src
	dst, err := run.StdoutPipe()

	err = run.Start()
	if err != nil {
		fmt.Println("RunStart Error", err)
		return err
	}
	defer func(dst io.ReadCloser) {
		err := dst.Close()
		if err != nil {

		}
	}(dst)

	defer func(Process *os.Process) {
		wait, err := Process.Wait()
		if err != nil {
			return
		}
		if !wait.Success() {
			err := Process.Kill()
			if err != nil {
				log.Println(err)
			}
			fmt.Println("killed", Process.Pid)
		}
	}(run.Process)

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	for {
		// read data from ffmpeg stdout
		pcm := make([]int16, frameSize*channels)
		err = binary.Read(dst, binary.LittleEndian, &pcm)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			fmt.Println("error reading from ffmpeg stdout", err)
			break
		}
		opus, err := opusEncoder.Encode(pcm, frameSize, maxBytes)
		if err != nil {
			fmt.Println("opus encode error", err)
		}
		chunks <- opus
	}

	done <- struct{}{}

	return nil
}
