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
		tpl, err := template.New("metrics").Parse(htmlTemplate)
		if err != nil {
			GlobalSugar.Fatalw(err.Error(), "event", "get index")
		}
		tpl.ExecuteTemplate(w, "metrics", Data.MetricsGauge)
		tpl.ExecuteTemplate(w, "metrics", Data.MetricsCounter)
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}
