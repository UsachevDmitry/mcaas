package internal

import (
	"html/template"
	"net/http"
)

const htmlTemplate = `{{ range $key, $value := . }}
<p>Key: {{ $key }}<br>Value: {{ $value }}</p>
{{ end }}`

func HandleIndex() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//ContentType := r.Header.Get("Content-Type")
		tpl, err := template.New("metrics").Parse(htmlTemplate)
		if err != nil {
			GlobalSugar.Fatalw(err.Error(), "event", "get index")
		}
		tpl.ExecuteTemplate(w, "metrics", Data.MetricsGauge)
		tpl.ExecuteTemplate(w, "metrics", Data.MetricsCounter)
		w.Header().Set("Accept", "text/html")
		w.Header().Set("Accept-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}
