package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"estiam/dictionary"

	"github.com/gorilla/mux"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	d := dictionary.New()

	stopServer := make(chan bool) // Canal pour arrêter le serveur Goby

	go runGobyServer(d, stopServer) // Démarrer le serveur Goby en arrière-plan

	runMode("cli", d, reader) // Lancer le mode CLI par défaut

	stopServer <- true // Arrêter le serveur Goby une fois que le mode CLI est terminé
}

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Println("Enter word:")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Println("Enter definition:")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.Add(word, definition)
	fmt.Println("Added.")
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Println("Enter word:")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	entry, err := d.Get(word)
	if err != nil {
		fmt.Println("Word not found.")
	} else {
		fmt.Println("Definition:", entry.String())
	}
}

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Println("Enter word:")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.Remove(word)
	fmt.Println("Removed.")
}

func actionList(d *dictionary.Dictionary) {
	words, entries := d.List()
	for _, word := range words {
		fmt.Println(word, ":", entries[word].String())
	}
}

func runMode(mode string, d *dictionary.Dictionary, reader *bufio.Reader) {
	for {
		fmt.Println("\nChoose an action [add/define/remove/list/exit]:")
		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(action)

		switch action {
		case "add":
			actionAdd(d, reader)
		case "define":
			actionDefine(d, reader)
		case "remove":
			actionRemove(d, reader)
		case "list":
			actionList(d)
		case "exit":
			if mode == "cli" {
				return
			}
		default:
			fmt.Println("Unknown action.")
		}
	}
}

func runGobyServer(d *dictionary.Dictionary, stopServer chan bool) {
	r := mux.NewRouter()

	r.HandleFunc("/define/{word}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]

		entry, err := d.Get(word)
		if err != nil {
			http.NotFound(w, r)
		} else {
			fmt.Fprintf(w, "Definition: %s", entry.String())
		}
	}).Methods("GET")

	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		word := r.FormValue("word")
		definition := r.FormValue("definition")
		d.Add(word, definition)
		fmt.Fprint(w, "Added.")
	}).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    "localhost:8080",
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	<-stopServer // Attendre la demande d'arrêt du serveur

	// Arrêter le serveur
	if err := srv.Close(); err != nil {
		log.Fatal(err)
	}
}
