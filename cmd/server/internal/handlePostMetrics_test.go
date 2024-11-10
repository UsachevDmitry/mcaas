package internal

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HandlePostMetrics(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		nameTest string
		dataType string
		name     string
		value    string
		request  string
		want     want
	}{
		{
			nameTest: "Test counter #1",
			request:  "/update/counter/t1/11",

			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			nameTest: "Test gauge #2",
			request:  "/update/gauge/t2/33.3",

			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			nameTest: "Test bad data type #3",
			request:  "/update/none/t1/11",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			nameTest: "Test empty value type #4",
			request:  "/update/counter/test6/",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {
			router := mux.NewRouter()
			router.HandleFunc("/update/{type}/{name}/{value}", WithLoggingHandlePostMetrics(HandlePostMetrics())).Methods(http.MethodPost)
			request := httptest.NewRequest(http.MethodPost, tt.request, http.NoBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
