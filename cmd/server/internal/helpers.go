package internal
import (
	"net/http"
	"time"
	"encoding/json"
	"strings"
	"compress/gzip"
	"io"
	"fmt"
	"bytes"
)

type Message struct {
	Message string `json:"message"`
}

var GlobalStatusCode int

func WithLoggingPost(h http.Handler) http.HandlerFunc { //func(w http.ResponseWriter, r *http.Request) {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		GlobalSugar.Infoln(
			"statusCode", GlobalStatusCode,
			"size", r.Header.Get("Content-Length"),
		)
	}
	return logFn
}

func WithLoggingGet(h http.Handler) http.HandlerFunc {//func(w http.ResponseWriter, r *http.Request) {
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

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
} 

func GzipHandle(h http.Handler) http.HandlerFunc {//func(w http.ResponseWriter, r *http.Request) {
	ArchFn := func(w http.ResponseWriter, r *http.Request) {
			// проверяем, что клиент поддерживает gzip-сжатие
			// это упрощённый пример. В реальном приложении следует проверять все
			// значения r.Header.Values("Accept-Encoding") и разбирать строку
			// на составные части, чтобы избежать неожиданных результатов

			if r.Header.Get("Content-Encoding") == "gzip" {
				r.Body = Decompress(r.Body)
			}

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				// если gzip не поддерживается, передаём управление
				// дальше без изменений
				h.ServeHTTP(w, r)
				return
			}

			// создаём gzip.Writer поверх текущего w
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			// передаём обработчику страницы переменную типа gzipWriter для вывода данных
			h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
	return ArchFn
}


func WriteHeaderAndSaveStatus(statusCode int, ContentType string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(statusCode)
	GlobalStatusCode = statusCode
}

func PostMetricAnswer(name string, dataType string, w http.ResponseWriter, r *http.Request){
	var CounterValueInt64 int64
	var GaugeValueFloat64 float64
	var ContentType string
	ContentType = r.Header.Get("Content-Type")

	if dataType == "counter" {
		CounterValue, exists := Data.GetCounter(name)
			if !exists {
				WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
				return
			}
		CounterValueInt64 = int64(CounterValue)
		var metrics = Metrics{
			ID: name,    
			MType: dataType,
			Delta: &CounterValueInt64,
		}
		requestBody, err := json.Marshal(metrics)
		if err != nil {
			GlobalSugar.Errorln("Error marshaling JSON:", err)
			return
		}
		//compressedJSONBody, _ := Compress(requestBody)
		w.Header().Set("Content-Type", ContentType)
		//w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(requestBody)
	} else if dataType == "gauge" {
		GaugeValue, exists := Data.GetGauge(name)
		if !exists {
			WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
			return
		}
		GaugeValueFloat64 = float64(GaugeValue)
		var metrics = Metrics{
			ID: name,    
			MType: dataType,
	        Value: &GaugeValueFloat64,
		}
		requestBody, err := json.Marshal(metrics)
		if err != nil {
			GlobalSugar.Errorln("Error marshaling JSON:", err)
			return
		}

		//compressedJSONBody, _ := Compress(requestBody)
		
		w.Header().Set("Content-Type", ContentType)
		//w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(requestBody)
	} else {
		WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
	}
}

// func Decompress(data []byte) ([]byte, error) {
//     // переменная r будет читать входящие данные и распаковывать их
//     r ,_ := gzip.NewReader(bytes.NewReader(data))
//     defer r.Close()

//     var b bytes.Buffer
//     // в переменную b записываются распакованные данные
//     _, err := b.ReadFrom(r)
//     if err != nil {
//         return nil, fmt.Errorf("failed decompress data: %v", err)
//     }

//     return b.Bytes(), nil
// } 

		// // Чтение сжатых данных из тела запроса io.ReadCloser r *http.Request
func Decompress(r io.ReadCloser) *gzip.Reader {
	body, err2 := io.ReadAll(r)
	if err2 != nil {
		fmt.Println("Error reading request body:", err2)
		return nil
	}
	// Распаковка данных
	//var data map[string]interface{}
	reader, err2 := gzip.NewReader(bytes.NewReader(body))
	if err2 != nil {
		fmt.Println("Error creating gzip reader:", err2)
		return nil
	}
	defer reader.Close()
	return reader
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	// Создаём переменную w — в неё будут записываться входящие данные,
	// которые будут сжиматься и сохраняться в bytes.Buffer
	w := gzip.NewWriter(&b)
	// Запись данных
	_ , err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	// Обязательно нужно вызвать метод Close() — в противном случае часть данных
	// может не записаться в буфер b; если нужно выгрузить все упакованные данные
	// в какой-то момент сжатия, используйте метод Flush()
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	// Переменная b содержит сжатые данные
	return b.Bytes(), nil
}
