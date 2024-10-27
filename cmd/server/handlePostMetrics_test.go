package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)
func Test_handlePostMetrics(t *testing.T) {
    type want struct {
        statusCode  int
    }
    tests := []struct {
        nameTest    string
		dataType string
		name string
		value string
		request string
        want    want
    }{
        {
            nameTest: "Test counter #1",
			request: "/update/counter/t1/11",

            want: want{
                statusCode:  http.StatusOK,
            },            
    	},
		{
            nameTest: "Test gauge #2",
			request: "/update/gauge/t2/33.3",

            want: want{
                statusCode:  http.StatusOK,
            },            
    	},
		{
            nameTest: "Test bad data type #3",
			request: "/update/none/t1/11",
            want: want{
                statusCode:  http.StatusBadRequest,
            },            
    	},
		{
            nameTest: "Test empty type type #4",
			request: "/update//testtest/11",
            want: want{
                statusCode:  http.StatusNotFound,
            },            
    	},
		{
            nameTest: "Test empty name type #5",
			request: "/update/counter//11",
            want: want{
                statusCode:  http.StatusNotFound,
            },            
    	},
		{
            nameTest: "Test empty value type #6",
			request: "/update/counter/test6/",
            want: want{
                statusCode:  http.StatusNotFound,
            },            
    	},
}

for _, tt := range tests {
	t.Run(tt.nameTest, func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, tt.request, nil)
		w := httptest.NewRecorder()
		handlePostMetrics(w, request)
		result := w.Result()
		assert.Equal(t, tt.want.statusCode, result.StatusCode)
	})
	}
}
