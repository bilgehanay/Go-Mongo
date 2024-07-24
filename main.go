package main

import (
	"log"
)

func main() {
	e.Go(func() error {
		return HttpServer.ListenAndServe()
	})

	if err := e.Wait(); err != nil {
		log.Fatal(err)
	}
}
