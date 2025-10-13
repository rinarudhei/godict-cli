package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var ErrNoAudio = errors.New("no audio available")

// playAudio plays the pronunciation audio for the given dictionary entry
func playAudio(entry DictionaryEntry) error {
	// Find first available audio URL
	audioURL := ""
	for _, p := range entry.Phonetics {
		if p.Audio != "" {
			audioURL = p.Audio
			break
		}
	}

	if audioURL == "" {
		return ErrNoAudio
	}

	tmpFile, err := downloadAudio(audioURL)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	return playAudioFile(tmpFile)
}

// downloadAudio downloads the audio file and returns the temp file path
func downloadAudio(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download audio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download audio: status %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "godict-*.mp3")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write audio: %w", err)
	}

	return tmpFile.Name(), nil
}

// playAudioFile plays an audio file using system commands
func playAudioFile(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", filePath)
	case "linux":
		if _, err := exec.LookPath("mpg123"); err == nil {
			cmd = exec.Command("mpg123", "-q", filePath)
		} else if _, err := exec.LookPath("mpv"); err == nil {
			cmd = exec.Command("mpv", "--really-quiet", filePath)
		} else if _, err := exec.LookPath("ffplay"); err == nil {
			cmd = exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", filePath)
		} else {
			return errors.New("no audio player found (install mpg123, mpv, or ffplay)")
		}
	case "windows":
		if _, err := exec.LookPath("mpv"); err == nil {
			cmd = exec.Command("mpv", "--really-quiet", filePath)
		} else {
			absPath, _ := filepath.Abs(filePath)
			psCmd := fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", absPath)
			cmd = exec.Command("powershell", "-Command", psCmd)
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}
