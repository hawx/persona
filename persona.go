package persona

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type personaResponse struct {
	Status string `json:"status"`
	Email  string `json:"email"`
}

func assert(audience, assertion string) (string, error) {
	params := url.Values{}
	params.Add("assertion", assertion)
	params.Add("audience", audience)

	resp, err := http.PostForm("https://verifier.login.persona.org/verify", params)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var f personaResponse
	err = json.Unmarshal(body, &f)

	if err != nil {
		return "", err
	}

	if f.Status != "okay" {
		return "", errors.New("Status not okay")
	}

	return f.Email, nil
}
