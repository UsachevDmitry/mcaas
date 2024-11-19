package internal

import (
	"net/http"
	"encoding/json"
)

func HandleGetMetricsJSON() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType = "application/json"
		var metrics Metrics
		// if r.Header.Get("Content-Encoding") == "gzip" {
		// 	r.Body = Decompress(r.Body)
		// }
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&metrics)
		if err != nil {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
			return
		}
		PostMetricAnswer(metrics.ID, metrics.MType, w)
	}
	return http.HandlerFunc(fn)
}