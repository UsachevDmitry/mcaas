package main

import (
	"github.com/UsachevDmitry/mcaas/cmd/server/internal"
	"github.com/gorilla/mux"
	"net/http"
	//"sync"
	"time"
	"os/signal"
	"os"
	"context"
	"syscall"
	"fmt"
)

func main() {
	//var wg sync.WaitGroup

	internal.GetConfig()
	internal.ImportDataFromFile(*internal.FileStoragePath, *internal.Restore)

	ctx1, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем функцию каждые 300 секунд
	ticker := time.NewTicker(300 * time.Second)
	go func() {
		for {
			select {
			case <-ctx1.Done():
				fmt.Println("Context canceled")
				return
			case <-ticker.C:
				internal.SaveDataInFile(time.Duration(*internal.StoreInterval), *internal.FileStoragePath)
			}
		}
	}()

	// Ждем завершения работы программы
	// fmt.Println("Программа запущена, ждем завершения...")
	// <-ctx1.Done()
	// fmt.Println("Программа завершена")

	// wg.Add(1)
	// go func() {
	// 	internal.SaveDataInFile(time.Duration(*internal.StoreInterval), *internal.FileStoragePath)
	// 	defer wg.Done()
	// }()

	router := mux.NewRouter()
	router.HandleFunc("/", internal.WithLoggingGet(internal.GzipHandle(internal.HandleIndex()))).Methods(http.MethodGet)
	router.HandleFunc("/update/", internal.WithLoggingGet(internal.GzipHandle(internal.HandlePostMetricsJSON()))).Methods(http.MethodPost)
	router.HandleFunc("/update/{type}/{name}/{value}", internal.WithLoggingGet(internal.GzipHandle(internal.HandlePostMetrics()))).Methods(http.MethodPost)
	router.HandleFunc("/value/", internal.WithLoggingGet(internal.GzipHandle(internal.HandleGetMetricsJSON()))).Methods(http.MethodPost)
	router.HandleFunc("/value/{type}/{name}", internal.WithLoggingGet(internal.GzipHandle(internal.HandleGetValue()))).Methods(http.MethodGet)

	internal.Logger()
	internal.GlobalSugar.Infow(
		"Starting server",
		"ADDRESS", *internal.Addr,
		"STORE_INTERVAL", *internal.StoreInterval,
		"FILE_STORAGE_PATH", *internal.FileStoragePath,
		"RESTORE", *internal.Restore,
	)

	srv := &http.Server{
		Addr: *internal.Addr,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		internal.GlobalSugar.Infoln("Graceful shutdown signal received")
		internal.SaveDataInFile(time.Duration(0), *internal.FileStoragePath)
		srv.Shutdown(ctx)
		os.Exit(0)
	}()

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		internal.GlobalSugar.Infoln("Error starting server:", err)
		os.Exit(1)
	}
}
