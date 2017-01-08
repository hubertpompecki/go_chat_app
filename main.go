package main

import (
    "log"
    "net/http"
    "text/template"
    "path/filepath"
    "sync"
    "flag"
    "github.com/hpompecki/trace"
    "os"
    "github.com/stretchr/gomniauth"
    "github.com/stretchr/objx"
    "github.com/stretchr/gomniauth/providers/google"
)

// templ represents a single template
type templateHandler struct {
    once sync.Once
    filename string
    templ *template.Template
}

// ServeHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t.once.Do(func () {
        t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
    })
    data := map[string]interface{} {
        "Host": r.Host,
    }
    if authCookie, err := r.Cookie("auth"); err == nil {
        data["UserData"] = objx.MustFromBase64(authCookie.Value)
    }
    t.templ.Execute(w, data)
}

func main() {
    var addr = flag.String("addr", ":8080", "The addr of the application.")
    flag.Parse()
    // set up gomniauth
    gomniauth.SetSecurityKey("a84293hf49n9ncfvbge87yoavbt4y8ayv47bv")
    gomniauth.WithProviders(
        google.New("459088933832-1s0647qk8c12dueicg14sjai0fknnts3.apps.googleusercontent.com",
                   "jsCTygv154JsFyRFkd9eBrkV",
                   "http://localhost:8080/auth/callback/google"),
    )
    r := newRoom()
    r.tracer = trace.New(os.Stdout)
    // root
    http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
    http.Handle("/login", &templateHandler{filename: "login.html"})
    http.HandleFunc("/logout", logoutHandler)
    http.HandleFunc("/auth/", loginHandler)
    http.Handle("/room", r)
    // start the chatroom
    go r.run()
    // start the web server
    log.Println("Starting web server on", *addr)
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}