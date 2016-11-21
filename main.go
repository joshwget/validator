package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rancher/os/config"
)

func main() {
	index, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(index))
	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprint(w, err)
			return
		}
		cloudConfig := r.FormValue("cc")

		result, err := config.Validate([]byte(cloudConfig))
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		if result.Valid() {
			fmt.Fprint(w, "Valid!")
		} else {
			for _, desc := range result.Errors() {
				fmt.Fprintf(w, "%s<br>", desc)
			}
		}
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}
