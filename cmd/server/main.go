package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"github.com/UsachevDmitry/mcaas/internal/server"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	var wg sync.WaitGroup
	internal.Logger()
	internal.GetConfig()

	db, err := internal.SelectStorage(internal.Config)
	if err != nil {
		internal.GlobalSugar.Errorln(err)
	}
	
	db.CreateTableGauge(context.Background())
	db.CreateTableCounter(context.Background())

	internal.ImportDataFromFile(*internal.FileStoragePath, *internal.Restore)

	wg.Add(1)
	if !internal.FlagUsePosgresSQL {
		go func() {
			internal.SaveDataInFile(time.Duration(*internal.StoreInterval), *internal.FileStoragePath)
			defer wg.Done()
		}()
	} else {
		wg.Done()
	}

	router := mux.NewRouter()
	router.HandleFunc("/", internal.WithLoggingGet(internal.GzipHandle(internal.HandleIndex()))).Methods(http.MethodGet)
	router.HandleFunc("/update/", internal.WithLoggingGet(internal.GzipHandle(internal.HandlePostMetricsJSON()))).Methods(http.MethodPost)
	router.HandleFunc("/update/{type}/{name}/{value}", internal.WithLoggingGet(internal.GzipHandle(internal.HandlePostMetrics()))).Methods(http.MethodPost)
	router.HandleFunc("/value/", internal.WithLoggingGet(internal.GzipHandle(internal.HandleGetMetricsJSON()))).Methods(http.MethodPost)
	router.HandleFunc("/value/{type}/{name}", internal.WithLoggingGet(internal.GzipHandle(internal.HandleGetValue()))).Methods(http.MethodGet)
	router.HandleFunc("/ping", internal.WithLoggingGet(internal.GzipHandle(internal.HandleGetPing()))).Methods(http.MethodGet)
	router.HandleFunc("/updates/", internal.WithLoggingGet(internal.GzipHandle(internal.HandlePostMetricsListJSON()))).Methods(http.MethodPost)

	internal.GlobalSugar.Infow(
		"Starting server",
		"ADDRESS", *internal.Addr,
		"STORE_INTERVAL", *internal.StoreInterval,
		"FILE_STORAGE_PATH", *internal.FileStoragePath,
		"RESTORE", *internal.Restore,
	)

	srv := &http.Server{
		Addr:    *internal.Addr,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	defer stop()

	go func() {
		<-ctx.Done()
		internal.GlobalSugar.Infoln("Graceful shutdown signal received")
		internal.SaveDataInFile(time.Duration(0), *internal.FileStoragePath)
		srv.Shutdown(ctx)
		os.Exit(0)
	}()

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		internal.GlobalSugar.Infoln("Error starting server:", err)
		os.Exit(1)
	}
	wg.Wait()
}
