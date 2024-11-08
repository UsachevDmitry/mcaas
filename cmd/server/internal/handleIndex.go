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
	// Создание логгера

        // функция Now() возвращает текущее время
        start := time.Now()

        // эндпоинт /ping
        uri := r.RequestURI
        // метод запроса
        method := r.Method

        // точка, где выполняется хендлер pingHandler 
        h.ServeHTTP(w, r) // обслуживание оригинального запроса
		
        // Since возвращает разницу во времени между start 
        // и моментом вызова Since. Таким образом можно посчитать
        // время выполнения запроса.
        duration := time.Since(start)

        // отправляем сведения о запросе в zap
        GlobalSugar.Infoln(
            "uri", uri,
            "method", method,
            "duration", duration,
        )
        //fmt.Println("TEST")
    }
    // возвращаем функционально расширенный хендлер
    return logFn
} 