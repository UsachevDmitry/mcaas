package internal

import (
	"net/http"
	"encoding/json"
)

func HandleGetMetricsJson() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType string = "application/json"
		var dataType string
		var name string
		var metrics Metrics

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&metrics)
		if err != nil {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
			message := Message{Message: "provided json file is invalid."}
			json.NewEncoder(w).Encode(message)
			return
		}
		
		dataType = metrics.MType
		name = metrics.ID

		WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
		PostMetricAnswer(name, dataType, w)
	}
	return http.HandlerFunc(fn)
}