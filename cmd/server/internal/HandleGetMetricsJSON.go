package internal

import (
	"net/http"
	"encoding/json"
)

func HandleGetMetricsJSON() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType string
		var metrics Metrics
		ContentType = r.Header.Get("Content-Type")
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