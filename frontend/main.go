package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

func setCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "token",
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
	Token string `json:"token"`
	Code  string `json:"code"`
	State string `json:"state"`
}

func newError(err error, w http.ResponseWriter) {
	log.Printf("Parse Error: " + err.Error())
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "call back error")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://localhost:8080/v1/auth/github/login")
	if err != nil {
		newError(err, w)
		return
	}
	log.Println(res)

	location := res.Header.Get("Location")
	res.Header.Get("Location")

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		newError(err, w)
		return
	}
	log.Println(string(body))
	callback := new(Callback)
	err = json.Unmarshal(body, &callback)
	if err != nil {
		newError(err, w)
		return
	}

	cookie := &http.Cookie{
		Name:  "token",
		Value: callback.Code,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, location, res.StatusCode)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		newError(err, w)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		newError(err, w)
		return
	}

	postUrl := "http://localhost:8080/v1/auth/github/callback"
	val := url.Values{}
	val.Set("code", r.Form.Get("code"))
	val.Set("state", r.Form.Get("state"))
	log.Println(r)
	req, err := http.NewRequest("POST", postUrl, strings.NewReader(val.Encode()))
	if err != nil {
		newError(err, w)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cookie.Value)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		newError(err, w)
		return
	}
	log.Println(res)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		newError(err, w)
		return
	}
	log.Println(string(body))

	tpl := template.Must(template.ParseFiles("templates/index.html"))
	//tpl.Execute(w, u)
	tpl.Execute(w, nil)
}
