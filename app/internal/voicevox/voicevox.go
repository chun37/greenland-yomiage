package voicevox

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/xerrors"
)

type VoiceVox struct {
	HTTPClient *http.Client
}

func New() *VoiceVox {
	return &VoiceVox{
		&http.Client{},
	}
}

func (v *VoiceVox) Generate(text string) ([]byte, error) {
	ctx := context.Background()

	audioQuery, err := v.getAudioQuery(ctx, text)
	if err != nil {
		return nil, xerrors.Errorf("failed to")
	}

	reader, err := v.getAudioBinary(ctx, audioQuery)
	if err != nil {
		return nil, xerrors.Errorf("failed to get audioBinary: %w", err)
	}

	return reader, nil
}

func (v *VoiceVox) getAudioBinary(ctx context.Context, audioQuery []byte) ([]byte, error) {
	synthesisURL, err := url.Parse("http://127.0.0.1:50021/synthesis")
	if err != nil {
		return nil, xerrors.Errorf("failed to parse synthesis URL: %w", err)
	}

	query := url.Values{}
	query.Add("speaker", "8")
	synthesisURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, synthesisURL.String(), bytes.NewReader(audioQuery))
	if err != nil {
		return nil, xerrors.Errorf("failed to create request object: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := v.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to run HTTPClient.Do: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("cannot close body:", err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("failed to read res.Body: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, xerrors.Errorf("failed to http request: %v", body)
	}

	return body, nil
}

func (v *VoiceVox) getAudioQuery(ctx context.Context, text string) ([]byte, error) {
	audioQueryURL, err := url.Parse("http://127.0.0.1:50021/audio_query")
	if err != nil {
		return nil, xerrors.Errorf("failed to parse audioQuery URL: %w", err)
	}

	query := url.Values{}
	query.Add("speaker", "8")
	query.Add("text", text)
	audioQueryURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, audioQueryURL.String(), nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to create request object: %w", err)
	}

	res, err := v.HTTPClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("failed to run HTTPClient.Do: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("cannot close body:", err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("failed to read res.Body: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, xerrors.Errorf("failed to http request: %v", body)
	}

	return body, nil
}
