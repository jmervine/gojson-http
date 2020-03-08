package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
    "log"

	"github.com/ChimeraCoder/gojson"
)

var (
	Listen      string
	Port        int
	Template    string
	Tmpl        *template.Template
	defaultJson = `{ "example": { "from": { "json": true } } }`
	mutty       = sync.Mutex{}
)

type Result struct {
	Json, Struct string
}

type Handler struct{}

func init() {
    log.SetFlags(0)
    log.SetPrefix("app=gojson-http")
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	defer r.Body.Close()

	Tmpl, err := template.ParseFiles(Template)
	if err != nil {
        log.Fatalf("at=ServeHTTP error=%v", err)
	}

    log.Printf("at=ServeHTTP method=%s path=%s user-agent=%s took=%v",
        r.Method, r.URL.Path, r.Header["User-Agent"], time.Since(begin))

	res := Result{
		Json: defaultJson,
	}

	if strings.HasSuffix(r.URL.Path, "json") {
		fmt.Fprintln(w, fmt.Sprintf(`{ "example": { "from": { "path": "%s" } } }`, r.URL.String()))
		return
	}

	var src string
	if r.Method == "POST" {
		val := r.PostFormValue("json")
		res.Json = val
	} else {
		src = r.URL.Query().Get("src")
		if src != "" {
			res.Json = src
		}
	}

	if strings.HasPrefix(res.Json, "http") {

		// redirect wth to src param, if res.Json is path, but src path doesn't exist
		if src == "" {
			http.Redirect(w, r, r.URL.Path+"?src="+strings.TrimSpace(res.Json), 301)
			return
		}

		// fetch res.Json
		resp, err := http.DefaultClient.Get(strings.TrimSpace(res.Json))
		if err != nil {
            log.Printf("at=ServeHTTP method=%s path=%s user-agent=%s took=%v",
                r.Method, r.URL.Path, r.Header["User-Agent"], time.Since(begin))
            log.Printf("at=ServeHTTP error=%v", err)
			res.Struct = fmt.Sprintf("JSON Parse Error: %v\n", err)
			Tmpl.Execute(w, nil)
			return
		}

		read, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
            log.Printf("at=ServeHTTP method=%s path=%s user-agent=%s took=%v",
                r.Method, r.URL.Path, r.Header["User-Agent"], time.Since(begin))
            log.Printf("at=ServeHTTP error=%v", err)
			res.Struct = fmt.Sprintf("JSON Fetch Error: %v\n", err)
		}
		res.Json = string(read)
	}

	if out, e := gojson.Generate(strings.NewReader(res.Json), gojson.ParseJson, "MyJsonName", "main", []string{"json"}, false, true); e == nil {
		res.Struct = string(out)
	} else {
        log.Printf("at=ServeHTTP method=%s path=%s user-agent=%s took=%v",
            r.Method, r.URL.Path, r.Header["User-Agent"], time.Since(begin))
        log.Printf("at=ServeHTTP error=%v", e)
		res.Struct = fmt.Sprintf("JSON Parse Error: %v\n", e)
	}
	Tmpl.Execute(w, res)
}

func main() {
	// reload tempalate on SIGHUP
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP)
	go reloadTemplate(sigc)

	flag.IntVar(&Port, "port", 8080, "startup port")
	flag.StringVar(&Listen, "listen", "localhost", "listen address")
	flag.StringVar(&Template, "template", "index.html", "display template")
	flag.Parse()

	handler := Handler{}

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", Listen, Port),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

    log.Printf("at=main address=%s", server.Addr)
    log.Fatalf("at=main error=%s", server.ListenAndServe())
}

func reloadTemplate(sigc chan os.Signal) {
	for _ = range sigc {
        log.Print("at=reloadTemplate message=\"reloading template\"")
		t, e := template.ParseFiles(Template)
		if e != nil {
            log.Printf("at=reloadTemplate error=%v", e)
		}
		mutty.Lock()
		Tmpl = t
		mutty.Unlock()
        log.Println("at=reloadTemplate message=\"reloading template\"")
	}
}
