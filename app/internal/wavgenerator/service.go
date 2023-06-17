package wavgenerator

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Service interface {
	Generate(s string) ([]byte, error)
}

type VoiceVox struct {
	HTTPClient *http.Client
}

func NewVoiceVox() Service {
	return &VoiceVox{
		&http.Client{},
	}
}

func (v *VoiceVox) Generate(text string) ([]byte, error) {
	ctx := context.Background()
	audioQuerty, err := v.getAudioQuery(ctx, text)
	if err != nil {
		return nil, err
	}
	reader, err := v.getAudioBinary(ctx, audioQuerty)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func (v *VoiceVox) getAudioBinary(ctx context.Context, audioQuery []byte) ([]byte, error) {
	synthesisURL, err := url.Parse("http://127.0.0.1:50021/synthesis")
	query := url.Values{}
	query.Add("speaker", "8")
	synthesisURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, synthesisURL.String(), bytes.NewReader(audioQuery))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := v.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("read error:", err)
		}
	}(res.Body)

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (v *VoiceVox) getAudioQuery(ctx context.Context, text string) ([]byte, error) {
	audioQueryURL, err := url.Parse("http://127.0.0.1:50021/audio_query")
	if err != nil {
		return nil, err
	}
	query := url.Values{}
	query.Add("speaker", "8")
	query.Add("text", text)
	audioQueryURL.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, audioQueryURL.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := v.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("read error:", err)
		}
	}(res.Body)

	audioQuery, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return audioQuery, nil
}
