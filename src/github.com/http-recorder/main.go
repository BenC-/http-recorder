package main

import (
	"flag"
	"fmt"
	nethttp "net/http"
	"os"
	"os/signal"

	"github.com/http-recorder/fifo"
	"github.com/http-recorder/http"
	"github.com/http-recorder/log"
)

var recorderPort string
var retrieverPort string

/*
 * Usage :
 * ./http-recorder -recorderPort XXXXX -retrieverPort XXXXX
 * Default values are
 *    - recorderPort 12345
 *    - retrieverPort  23456
 */
func main() {

	fmt.Println(" _     _   _                                      _           \n" +
		"| |   | | | |                                    | |          \n" +
		"| |__ | |_| |_ _ __    _ __ ___  ___ ___  _ __ __| | ___ _ __ \n" +
		"| '_ \\| __| __| '_ \\  | '__/ _ \\/ __/ _ \\| '__/ _` |/ _ \\ '__|\n" +
		"| | | | |_| |_| |_) | | | |  __/ (_| (_) | | | (_| |  __/ |   \n" +
		"|_| |_|\\__|\\__| .__/  |_|  \\___|\\___\\___/|_|  \\__,_|\\___|_|   \n" +
		"              | |                                             \n" +
		"              |_|                                             \n")

	fmt.Println("starting http recorder...")
	flag.StringVar(&recorderPort, "recorderPort", "12345", "Port on which requests are catched and stored")
	flag.StringVar(&retrieverPort, "retrieverPort", "23456", "Port on which requests can be retrieved")
	flag.Parse()
	fifo.Init()
	go nethttp.ListenAndServe(fmt.Sprint(":", recorderPort), nethttp.HandlerFunc(http.RecorderHandler))
	go nethttp.ListenAndServe(fmt.Sprint(":", retrieverPort), nethttp.HandlerFunc(http.RetrieverHandler))

	log.RecorderInfo("Recorder is listening on port", recorderPort)
	log.RetrieverInfo("Retriever is listening on port", retrieverPort)

	waitForStop()
}

func waitForStop() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("interrup signal received (" + s.String() + "), shutting down server")
}
