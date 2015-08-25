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

	"github.com/jmervine/gojson-http/Godeps/_workspace/src/github.com/ChimeraCoder/gojson"
	"github.com/jmervine/gojson-http/Godeps/_workspace/src/gopkg.in/jmervine/readable.v1"
)

var (
	Listen      string
	Port        int
	Template    string
	Tmpl        *template.Template
	defaultJson = `{ "example": { "from": { "json": true } } }`
	mutty       = sync.Mutex{}
)
var log = readable.New().
	WithPrefix("http-gojson").
	WithOutput(os.Stdout).
	WithFlags(0)

type Result struct {
	Json, Struct string
}

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	defer r.Body.Close()

	Tmpl, err := template.ParseFiles(Template)
	if err != nil {
		log.Panic("at", "ServerHTTP", "error", err)
	}

	log.Log("at", "ServeHTTP", "method", r.Method, "path", r.URL.Path, "user-agent", r.Header["User-Agent"], "took", time.Since(begin))

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
			logError(r, begin, err)
			res.Struct = fmt.Sprintf("JSON Parse Error: %v\n", err)
			Tmpl.Execute(w, nil)
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
	Tmpl.Execute(w, res)
}

func logError(r *http.Request, t time.Time, e error) {
	log.Log("at", "logError", "method", r.Method, "path", r.URL.Path, "user-agent", r.Header["User-Agent"], "took", time.Since(t))
	log.Log("at", "logError", "error", e)
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

	log.Log("at", "main", "address", server.Addr)
	log.Fatal("at", "main", "error", server.ListenAndServe())
}

func reloadTemplate(sigc chan os.Signal) {
	for _ = range sigc {
		log.Log("at", "reloadTemplate", "message", "reloading template")
		t, e := template.ParseFiles(Template)
		if e != nil {
			log.Log("at", "reloadTemplate", "error", e)
		}
		mutty.Lock()
		Tmpl = t
		mutty.Unlock()
		log.Log("at", "reloadTemplate", "message", "template reloaded")
	}
}
