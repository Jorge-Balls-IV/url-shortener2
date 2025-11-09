package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/logging"
	"url-shortener/internal/response"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
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
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			alias:     "test_alias",
			url:       "",
			respError: "field URL is a required field",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Invalid URL",
			alias:     "test_alias",
			url:       "fdsfasdfasdf",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to save URL",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, test := range cases {

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if test.respError == "" || test.mockError != nil {
				// Говорим моку, что от указанной функции ожидаем, что будет передана канкретная url-строка и в качестве alias'a - любая строка. Описываем, что должно вернуться, и указываем, что это должно исполниться один раз
				urlSaverMock.On("SaveURL", test.url, mock.AnythingOfType("string")).
					Return(int64(1), test.mockError).Once()
			}

			handler := save.New(logging.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, test.url, test.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			respRec := httptest.NewRecorder() //Создаём тестовый рекордер чтобы записать туда данные

			handler.ServeHTTP(respRec, req)

			require.Equal(t, respRec.Code, http.StatusOK)

			body := respRec.Body.String()

			var resp response.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, test.respError, resp.Error)
		})
	}

}
