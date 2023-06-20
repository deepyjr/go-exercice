package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "estiam/dictionary"
)

func main() {
    d := dictionary.New()
    reader := bufio.NewReader(os.Stdin)

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
            return
        default:
            fmt.Println("Unknown action.")
        }
    }
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