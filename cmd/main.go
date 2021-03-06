package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/numeroSette/SRE-TEST-7/cmd/get-random-number-native/register"
	_ "github.com/numeroSette/SRE-TEST-7/cmd/get-random-number/register"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/numeroSette/SRE-TEST-7/internal/config"
	"github.com/numeroSette/SRE-TEST-7/internal/router"
)

func init() {

	config.
		Add(
			"http-service-listen-address",
			"HTTP_SERVICE_LISTEN_ADDRESS",
			string("0.0.0.0:8080"),
			"IP:PORT address to listen as service endpoint",
		).
		Add(
			"http-metrics-listen-address",
			"HTTP_METRICS_LISTEN_ADDRESS",
			string("0.0.0.0:8081"),
			"IP:PORT address to listen as metrics endpoint",
		)

}

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	config.Register().Load()

	router.Router.Use(router.PrometheusMiddleware)

	serviceServer := http.NewServeMux()
	serviceServer.Handle("/", router.Router)

	go func() {
		address := config.Get("http-service-listen-address").GetStringVal()

		log.Printf("service server started on: http://%s\n", address)
		if err := http.ListenAndServe(address, serviceServer); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	go func() {
		address := config.Get("http-metrics-listen-address").GetStringVal()

		log.Printf("metrics server started on: http://%s\n", address)
		if err := http.ListenAndServe(address, metricsServer); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("received signal: %v\n", sig)
		done <- true
	}()

	log.Println("awaiting signal")
	<-done
	log.Println("exiting")

}
