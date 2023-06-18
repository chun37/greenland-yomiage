package voicevox

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/xerrors"
)

func (v *VoiceVox) Add(word, pronunciation string, accent int) error {
	userDictURL, err := url.Parse("http://127.0.0.1:50021/user_dict_word")
	if err != nil {
		return xerrors.Errorf("failed to parse user dictionary URL: %w", err)
	}

	query := url.Values{}
	query.Add("surface", word)
	query.Add("pronunciation", pronunciation)
	query.Add("accent_type", strconv.Itoa(accent))
	userDictURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, userDictURL.String(), nil)
	if err != nil {
		return xerrors.Errorf("failed to create request object: %w", err)
	}

	res, err := v.HTTPClient.Do(req)
	if err != nil {
		return xerrors.Errorf("failed to run HTTPClient.Do: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("cannot close body:", err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return xerrors.Errorf("failed to read res.Body: %w", err)
	}

	if res.StatusCode != 200 {
		return xerrors.Errorf("failed to http request: %v", string(body))
	}

	return nil
}
