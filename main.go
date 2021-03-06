package main

import (
    "fmt"
    "net/http"
    "log"
    "time"
    "os"
    // "crypto/sha256"
    // "encoding/hex"

    "github.com/googollee/go-socket.io"
    "github.com/gorilla/mux"
)


// middleware for adding a cookie
func cookieMiddleware(h http.Handler) http.Handler  {
    // check if has cookie, else set cookie

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Some code before request...
        _, cookieFailed := r.Cookie("yummy")

        if cookieFailed != nil {
            
            expire := time.Now().AddDate(0, 0, 1)
            val := genId()

            cookie := http.Cookie{Name: "yummy", Value: val, Path: "/", Expires: expire, MaxAge: 86400}
            http.SetCookie(w, &cookie)
        }

        h.ServeHTTP(w, r)
        // Some code after request...
    })
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World! :O")
}


func main() {
    var port = os.Getenv("PORT")
    if port == "" {
        port = "8000"
        fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
    }

    port = ":" + port

    story := LoadStory()
    gm := newGameMaster(&story)
    // derp := Derp{}
    log.Println("Loaded Story:", story)

    // fmt.Println(story.getPath(1))
    // fmt.Println(story.hasEnded(7))

    // start the mux router
    r := mux.NewRouter()

    // serve static files
    s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
    r.PathPrefix("/static/").Handler(s)

    // handle the socket
    server, err := socketio.NewServer(nil)
    if err != nil {
        log.Fatal(err)
    }
    socketHandler := metaSocketHandler(&gm)
    server.On("connection", socketHandler)
    server.On("error", socketErrorHandler)
    r.Handle("/socket.io/", server)

    // main route
    r.Handle("/", cookieMiddleware(http.FileServer(http.Dir("./static/")) ))

    // API routes
    r.HandleFunc("/hello", helloHandler)

    log.Println("Started Server!")
    http.ListenAndServe(port, r)
}
