package internal

import (
	"html/template"
	"net/http"
	"time"
)

const htmlTemplate = `{{ range $key, $value := . }}
<p>Ключ: {{ $key }}<br>Значение: {{ $value }}</p>
{{ end }}`

func HandleIndex() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.New("metrics").Parse(htmlTemplate)
		if err != nil {
			GlobalSugar.Fatalw(err.Error(), "event", "get index")
		}
		tpl.ExecuteTemplate(w, "metrics", Data.MetricsGauge)
		tpl.ExecuteTemplate(w, "metrics", Data.MetricsCounter)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

func WithLoggingHandleIndex(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r)
		duration := time.Since(start)
		GlobalSugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	}
	return logFn
}
