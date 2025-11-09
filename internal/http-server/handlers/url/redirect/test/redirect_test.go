package test

import (
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"url-shortener/internal/logging"
	"url-shortener/internal/redirectCheck"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "http://google.com/",
		},
		/**{
			name:      "Empty_alias",
			alias:     "",
			url:       "http://google.com/",
			respError: "cannot redirect on empty alias",
		},
		{
			name:      "Url_Not_Found",
			alias:     "test_alias",
			url:       "",
			respError: "url not found",
			mockError: fmt.Errorf("unexpected error"),
		},**/
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if test.respError == "" || test.mockError != nil {
				urlGetterMock.On("GetURL", test.alias).
					Return(test.url, test.mockError).Once()
			}

			r := chi.NewRouter()

			r.Get("/{alias}",
				redirect.New(logging.NewDiscardLogger(),
					urlGetterMock))

			tServer := httptest.NewServer(r)
			defer tServer.Close()

			redirect, err := redirectCheck.GetRedirect(tServer.URL + "/" + test.alias)
			require.NoError(t, err)
			// Проверяем финальный url после редиректа
			require.Equal(t, test.url, redirect)
		})
	}
}
