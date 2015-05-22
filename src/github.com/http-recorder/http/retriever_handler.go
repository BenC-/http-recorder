package http

import (
	"encoding/json"
	"fmt"
	"github.com/http-recorder/entities"
	"github.com/http-recorder/fifo"
	nethttp "net/http"
	"strconv"
	"time"
)

func RetrieverHandler(w nethttp.ResponseWriter, r *nethttp.Request) {
	fmt.Println("[HTTP-RETRIEVER] (!) new client connection from", r.RemoteAddr)

	awaitingTime, headerNotAvailableOrMalformed := strconv.Atoi(r.Header.Get("Request-Timeout"))
	if headerNotAvailableOrMalformed != nil {
		awaitingTime = 3
	}
	awaitingChan := make(chan *entities.HttpRequest)
	stopChan := make(chan bool)
	go func() {
		if len(r.URL.Query()) != 0 { // Search by path
			fmt.Println("[HTTP-RETRIEVER] client asks for a specific request", r.URL.Query())

			var request *entities.HttpRequest
			var err error

			for key, value := range r.URL.Query() {
				request, err = fifo.FindBy(key, value[0])
				break
			}
			// Handle long polling
			for err != nil {
				select {
				case <-stopChan:
					return
				default:
					for key, value := range r.URL.Query() {
						request, err = fifo.FindBy(key, value[0])
						break
					}
					time.Sleep(1 * time.Second)
				}
			}
			awaitingChan <- request
		} else { // No search
			fmt.Println("[HTTP-RETRIEVER] client asks for any type of request")
			request, err := fifo.GetOldest()
			for err != nil {
				select {
				case <-stopChan:
					return
				default:
					request, err = fifo.GetOldest()
					time.Sleep(1 * time.Second)
				}
			}
			awaitingChan <- request
		}
	}()

	select {
	case <-time.After(time.Second * time.Duration(awaitingTime)):
		fmt.Println("[HTTP-RETRIEVER] sorry timeout reached, query returned no result, goodbye")
		w.WriteHeader(nethttp.StatusNotFound)
		stopChan <- true
	case request := <-awaitingChan:
		fmt.Println("[HTTP-RETRIEVER] return following request to client", request)
		json.NewEncoder(w).Encode(request)
	}

}
