package main

import (
    "os"
    "sync"
    "log"
    "fmt"
    "flag"
    "time"
    "strings"
    "syscall"
    "net/http"
    "os/signal"
    "io/ioutil"
    "html/template"

    "github.com/jmervine/gojson"
)

var (
    Listen string
    Port int
    Template, err = template.ParseFiles("./index.html");
    defaultJson = `{ "example": { "from": { "json": true } } }`
    mutty = sync.Mutex{}
)

type Result struct {
    Json, Struct string
}

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    begin := time.Now()

    defer r.Body.Close()

    log.Printf("=> %v %v 200 %v %v %s\n", r.Method, r.URL, r.Proto, r.Header["User-Agent"], fmt.Sprintf("%s", time.Since(begin)))

    res := Result{
        Json: defaultJson,
    }

    if strings.HasSuffix(r.URL.Path, "json") {
        fmt.Fprintln(w, `{ "example": { "from": { "path": true } } }`)
        return
    }

    if r.Method == "POST" {
        val := r.PostFormValue("json")
        res.Json = val
    }

    if strings.HasPrefix(res.Json, "http") {
        resp, err := http.DefaultClient.Get(strings.TrimSpace(res.Json))
        if err != nil {
            logError(r, begin, err)
            res.Struct = fmt.Sprintf("JSON Parse Error: %v\n", err)
            Template.Execute(w, nil)
            return
        }

        read, err := ioutil.ReadAll(resp.Body)
        resp.Body.Close()
        if err != nil {
            logError(r, begin, err)
            res.Struct = fmt.Sprintf("JSON Fetch Error: %v\n", err)
        }
        res.Json = string(read)
    }

    if out, e := json2struct.Generate(strings.NewReader(res.Json), "MyJsonName", "main"); e == nil {
        res.Struct = string(out)
    } else {
        logError(r, begin, err)
        res.Struct = fmt.Sprintf("JSON Parse Error: %v\n", err)
    }
    Template.Execute(w, res)
}

func logError(r *http.Request, t time.Time, e error) {
    log.Printf("=> %v %v 500 %v %v %s\n", r.Method, r.URL, r.Proto, r.Header["User-Agent"], fmt.Sprintf("%s", time.Since(t)))
    log.Printf("ERROR:\n%v\n\n", e)
}

func main() {

    if err != nil {
        log.Panic(err)
    }

    // reload tempalate on SIGHUP
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGHUP)
    go reloadTemplate(sigc)

    flag.IntVar(&Port, "port", 8080, "startup port")
    flag.StringVar(&Listen, "listen", "localhost", "listen address")
    flag.Parse()

    handler := Handler{}

    server := &http.Server{
        Addr:           fmt.Sprintf("%s:%d", Listen, Port),
        Handler:        handler,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    log.Printf("Starting at %v\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}

func reloadTemplate(sigc chan os.Signal) {
    for _ = range sigc {
        log.Print("Attempting to reload template!")
        t, e := template.ParseFiles("./index.html");
        if e == nil {
            mutty.Lock()
            Template = t
            mutty.Unlock()
            log.Print("Template reloaded!")
        }
    }
}

