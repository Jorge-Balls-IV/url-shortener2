package redirectCheck

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
)

// Возвращаем последний URL после редиректа
func GetRedirect(url string) (string, error) {
	const origin = "internal.redirectCheck.GetRedirect"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Останавливаемся после первого редиректа
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s: %w: %d", origin, ErrInvalidStatusCode, resp.StatusCode)

	}

	defer func(){
		resp.Body.Close()
	}()

	return resp.Header.Get("Location"), nil
}
