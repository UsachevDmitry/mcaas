package internal

import (
	"html/template"
	"log"
	"net/http"
)

const htmlTemplate = `{{ range $key, $value := . }}
<p>Ключ: {{ $key }}<br>Значение: {{ $value }}</p>
{{ end }}`

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("metrics").Parse(htmlTemplate)
	if err != nil {
		log.Fatalln(err)
	}
	tpl.ExecuteTemplate(w, "metrics", Data.MetricsGauge)
	tpl.ExecuteTemplate(w, "metrics", Data.MetricsCounter)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
}
