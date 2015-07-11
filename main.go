package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jmervine/gojson-http/Godeps/_workspace/src/github.com/ChimeraCoder/gojson"
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

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()

	defer r.Body.Close()

	Tmpl, err := template.ParseFiles(Template)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("=> %v %v 200 %v %v %s\n", r.Method, r.URL, r.Proto, r.Header["User-Agent"], fmt.Sprintf("%s", time.Since(begin)))

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
	log.Printf("=> %v %v 500 %v %v %s\n", r.Method, r.URL, r.Proto, r.Header["User-Agent"], fmt.Sprintf("%s", time.Since(t)))
	log.Printf("ERROR:\n%v\n\n", e)
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

	log.Printf("Starting at %v\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func reloadTemplate(sigc chan os.Signal) {
	for _ = range sigc {
		log.Print("Attempting to reload template!")
		t, e := template.ParseFiles(Template)
		if e == nil {
			mutty.Lock()
			Tmpl = t
			mutty.Unlock()
			log.Print("Template reloaded!")
		}
	}
}
