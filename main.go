package main

import (
    "log"
    "net/http"
    "chat_app/server"
)

func main() {
    hub := server.NewHub()
    go hub.Run()

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        server.ServeWs(hub, w, r)
    })

    fs := http.FileServer(http.Dir("./client"))
    http.Handle("/", fs)

    log.Println("Server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
