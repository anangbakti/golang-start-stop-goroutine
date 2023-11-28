package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

// DB is the global variable for the database connection
var (
	nmr      chan int
	stopOnce sync.Once
)

func main() {
	http.HandleFunc("/task5/", task5)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func task5(w http.ResponseWriter, r *http.Request) {
	goIgnition := r.URL.Path[len("/task5/"):]

	switch goIgnition {
	case "start":
		workerStart()
		w.Write([]byte("Your worker server started\n"))
	case "stop":
		workerStop()
		w.Write([]byte("Your worker server stoped\n"))
	case "exit":
		os.Exit(3)
	default:
		http.Error(w, "Invalid command text. Use 'start' or 'stop' to command go concurrency. And 'exit' to exit program.", http.StatusBadRequest)
		return
	}
}

func workerStart() {
	stopOnce = sync.Once{} // Reset the stopOnce sync.Once
	nmr = make(chan int)
	go func() {
		fmt.Println("worker start ...")
		for {
			select {
			case <-time.After(2 * time.Second):
				fmt.Println(rand.Intn(100))
			case <-nmr:
				return
			}
		}
	}()
}

func workerStop() {
	fmt.Println("worker done")
	stopOnce.Do(func() {
		if nmr != nil {
			close(nmr)
		}
	})
}

// Generator returns a channel that produces the numbers 1, 2, 3,…
// To stop the underlying goroutine, send a number on this channel.
func Generator() chan int {
	go func() {
		n := 1
		for {
			select {
			case nmr <- n:
				n++
			case <-nmr:
				return
			}
		}
	}()
	return nmr
}
