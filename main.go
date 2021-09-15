package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "Hello world")
	case http.MethodPost:
		if err := req.ParseForm(); err != nil {
			fmt.Println("Something went bad")
			fmt.Fprintln(w, "Something went bad")
			return
		}
		for key, value := range req.PostForm {
			fmt.Println(key, "=>", value)
		}
		fmt.Fprintf(w, "Information received: %v\n", req.PostForm)
	}
}

func dateHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		currentTime := time.Now()
		fmt.Fprintf(w, "%s", currentTime.Format("03h04"))
	}
}

func addHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		if err := req.ParseForm(); err != nil { // Parsing des paramètres envoyés
			fmt.Println("Something went bad") // par le client et gestion d’erreurs
			fmt.Fprintln(w, "Something went bad")
			return
		}
		for key, value := range req.PostForm { // On print les clés et valeurs des
			fmt.Println(key, "=>", value) // données envoyés par le clients
		}
		fmt.Fprintf(w, "Information received: %v\n", req.PostForm)
		// Sauvegarde des données reçues
		saveFile, err := os.OpenFile("./save.data", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		defer saveFile.Close()

		w := bufio.NewWriter(saveFile)
		if err == nil {
			fmt.Fprintf(w, "%v:%v:\n", req.PostForm["entries"][0], req.PostForm["author"][0])
		}
		w.Flush()

	}
}

func entriesHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		filerc, err := os.Open("./save.data")
		if err != nil {
			log.Fatal(err)
		}
		defer filerc.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(filerc)
		contents := buf.String()

		split := strings.Split(contents, ":")

		for k := range split {
			if k%2 == 0 {
				fmt.Fprintf(w, split[k])
			}
		}
	}
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", dateHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/entries", entriesHandler)
	http.ListenAndServe(":4567", nil)
}
