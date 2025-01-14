package lib

import (
	"fmt"
	"io"
	"net/http"
)

type ICredentials interface {
	GetCredentials() (string, string)
}

func HttpRequest(method string, url string, credentials ICredentials) ([]byte, error) {
	var err error

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unexpected error creating an http request %w", err)
	}

	username, password := credentials.GetCredentials()
	if username != "" && password != "" {
		request.SetBasicAuth(username, password)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to send http request %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("error calling http url \"%v\". status code: %v", url, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unexpected error reading http response body %w", err)
	}

	return body, nil
}
