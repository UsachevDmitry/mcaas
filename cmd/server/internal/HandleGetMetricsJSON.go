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
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&metrics)
		if err != nil {
			WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
			return
		}
		PostMetricAnswer(metrics.ID, metrics.MType, w, r)
	}
	return http.HandlerFunc(fn)
}