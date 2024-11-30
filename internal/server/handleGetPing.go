package internal

import (
	"context"
	"net/http"
	"time"
)

func HandleGetPing() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := DB.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			GlobalSugar.Panicln(err)
		}

		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}
