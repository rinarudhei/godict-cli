package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type DictionaryEntry struct {
	Word      string     `json:"word"`
	Phonetic  string     `json:"phonetic"`
	Phonetics []Phonetic `json:"phonetics"`
	Origin    string     `json:"origin"`
	Meanings  []Meaning  `json:"meanings"`
}

type Phonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio,omitempty"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type Definition struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example,omitempty"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
}

var (
	ErrConnection      = errors.New("connection error")
	ErrInvalidResponse = errors.New("invalid response")
	ErrNotFound        = errors.New("data not found")
)

func newHTTPClient() *http.Client {
	timeout := viper.GetDuration("api-timeout")
	return &http.Client{
		Timeout: time.Second * timeout,
	}
}

func lookDictionary(apiURL, word string) ([]DictionaryEntry, error) {
	var resp []DictionaryEntry
	r, err := newHTTPClient().Get(apiURL + word)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if r.StatusCode != http.StatusOK {
		return nil, ErrInvalidResponse
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidResponse, err)
	}

	return resp, nil
}
