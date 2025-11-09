package test

import (
	"net/http"
	"net/url"
	"path"
	"testing"
	"url-shortener2/internal/random"
	"url-shortener2/internal/redirectCheck"
	"url-shortener2/internal/response"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {

	//Формируем базовый url, к которому будет образаться клиент
	url := url.URL{
		Scheme: "http",
		Host:   host,
	}

	// Создаём специальный клиент, с помощью которого будут отправляться запросы
	expect := httpexpect.Default(t, url.String())

	expect.POST("/url").WithJSON(
		response.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("Jorge_Balls_IV", "84985538456").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")

}

func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   random.NewRandomString(12),
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			url := url.URL{
				Scheme: "http",
				Host:   host,
			}

			expect := httpexpect.Default(t, url.String())

			//save
			resp := expect.POST("/url").
				WithJSON(response.Request{
					URL:   test.url,
					Alias: test.alias,
				}).
				WithBasicAuth("Jorge_Balls_IV", "84985538456").
				Expect().
				Status(http.StatusOK).
				JSON().Object()

			if test.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(test.error)

				return
			}

			alias := test.alias

			if alias != "" {
				resp.Value("alias").String().IsEqual(alias)
			} else {
				resp.Value("alias").String().NotEmpty()
				alias = resp.Value("alias").String().Raw()
			}

			//redirect

			testRedirect(t, alias, test.url)

			//remove

			respDel := expect.DELETE("/"+path.Join("url", "remove", alias)).
				WithBasicAuth("Jorge_Balls_IV", "84985538456").
				Expect().Status(http.StatusOK).
				JSON().Object()

			respDel.Value("status").String().IsEqual("OK")

			//redirect again

			testNoRedirect(t, alias)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	url := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirURL, err := redirectCheck.GetRedirect(url.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirURL)
}

func testNoRedirect(t *testing.T, alias string) {
	url := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	_, err := redirectCheck.GetRedirect(url.String())

	require.Error(t, err, redirectCheck.ErrInvalidStatusCode)
}
