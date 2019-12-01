package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
)

func setCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: createRand(),
	}
	http.SetCookie(w, cookie)

	tpl := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
	tpl.Execute(w, cookie)
}

const (
	letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	indexBit  = 6
	indexMask = 1<<indexBit - 1
	indexMax  = 63 / indexBit
)

func createRand() (randVal string) {
	randSource := rand.NewSource(time.Now().UnixNano())
	n := 32
	b := make([]byte, n)
	cache, remain := randSource.Int63(), indexMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSource.Int63(), indexMax
		}
		index := int(cache & indexMask)
		if index < len(letters) {
			b[i] = letters[index]
			i--
		}
		cache >>= indexBit
		remain--
	}
	randVal = string(b)
	return
}

func main() {
	http.HandleFunc("/", setCookie)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalln(err)
	}
}

// Callback is callback from github
type Callback struct {
	Session string `json:"session_id"`
	Code    string `json:"code"`
	State   string `json:"state"`
}

func newError(err error, w http.ResponseWriter) {
	log.Printf("Parse Error: " + err.Error())
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "call back error")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		newError(err, w)
		return
	}

	val, err := json.Marshal(Callback{
		Session: cookie.Value,
	})
	if err != nil {
		newError(err, w)
		return
	}
	res, err := http.Post("http://localhost:8080/v1/auth/github/login", "application/json", bytes.NewBuffer(val))
	if err != nil {
		newError(err, w)
		return
	}
	log.Println(res)

	location := res.Header.Get("Location")
	res.Header.Get("Location")

	http.Redirect(w, r, location, res.StatusCode)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		newError(err, w)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		newError(err, w)
		return
	}

	val, err := json.Marshal(Callback{
		Session: cookie.Value,
		Code:    r.Form.Get("code"),
		State:   r.Form.Get("state"),
	})
	if err != nil {
		newError(err, w)
		return
	}
	res, err := http.Post("htp://localhost:8080/v1/auth/github/callback", "application/json", bytes.NewBuffer(val))
	if err != nil {
		newError(err, w)
		return
	}
	log.Println(res)

	//defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		newError(err, w)
		return
	}

	/*
		src := p.Config.TokenSource(ctx, p.Token)
		httpClient := oauth2.NewClient(ctx, src)

		u, err := p.GetUsers(ctx, httpClient)
		if err != nil {
			log.Printf("get user error:" + err.Error())
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "get user error")
			return
		}
	*/
	tpl := template.Must(template.ParseFiles("templates/index.html"))
	//tpl.Execute(w, u)
	tpl.Execute(w, nil)
}
