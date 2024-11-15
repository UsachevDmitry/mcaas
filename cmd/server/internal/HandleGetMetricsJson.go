package internal

import (
	"net/http"
	//"encoding/json"
)

func HandleGetMetricsJson() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType string = "application/json"

		//json.NewEncoder(w).Encode(DB)
		// GlobalSugar.Infoln("Get")
		WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
	}
	return http.HandlerFunc(fn)
}