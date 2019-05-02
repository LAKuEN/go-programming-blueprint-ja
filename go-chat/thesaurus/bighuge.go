package thesaurus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type BigHuge struct {
	APIKey  string
	BaseURL string
}

func NewBigHuge(apiKey string) BigHuge {
	baseURL := "http://words.bighugelabs.com/api/2/<apikey>/<word>/json"
	return BigHuge{APIKey: apiKey, BaseURL: baseURL}
}

type synonyms struct {
	Noun *words `json:"noun"`
	Verb *words `json:"verb"`
}

type words struct {
	Syn []string `json:"syn"`
}

func (b *BigHuge) Synonyms(term string) ([]string, error) {
	var syns []string
	url := strings.Replace(b.BaseURL, "<word>", term, 1)
	url = strings.Replace(url, "<apikey>", b.APIKey, 1)
	response, err := http.Get(url)
	if err != nil {
		return syns, fmt.Errorf("bighuge: failed to find synonyms of %s: %v", term, err)
	}

	var data synonyms
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return syns, err
	}

	if data.Noun.Syn != nil {
		syns = append(syns, data.Noun.Syn...)
	}
	if data.Verb.Syn != nil {
		syns = append(syns, data.Verb.Syn...)
	}

	return syns, nil
}
