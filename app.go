package main

import (
	//"io"
	//"io/ioutil"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var chatLog bytes.Buffer

//var chatLog make([]byte,2)

// 	{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func tempUser() string {

	r := make([]byte, 8)
	_, err := rand.Read(r)

	if err != nil {
		log.Println(err)
	}
	log.Println(r)
	rStr := base64.URLEncoding.EncodeToString(r)

	return rStr
}

func receiveHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		lang := r.FormValue("lang")
		log.Println("Received: ", lang)
		chatLog.WriteString("<p>guest013: " + lang + "</p>")
		//chatLog.WriteString("<hr>")
		w.Write(chatLog.Bytes())
	}
}

func getComments(w http.ResponseWriter, r *http.Request) {
	w.Write(chatLog.Bytes())
}

func intro(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("./index.html")

	if err != nil {
		http.NotFound(w, r)
	} else {

		data := struct {
			Title string
			Items []string
		}{
			Title: "A test site",
			Items: []string{
				"My photos",
				"My blog",
			},
		}

		err = t.Execute(w, data)
		checkErr(err)
	}

}

var validPath = regexp.MustCompile("^/(receive|getComments|view)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

// main function starts the HTTP server
func main() {

	chatLog.WriteString("<hr><p><b>COMMENTS</b></p>")
	chatLog.WriteString("<p>rand0m: Maecenas nec odio et ante tincidunt tempus. Donec vitae sapien ut libero venenatis faucibus. Nullam quis ante.</p>")

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", intro)
	http.HandleFunc("/receive", makeHandler(receiveHandler))
	http.HandleFunc("/getComments", makeHandler(getComments))

	log.Printf("Server running at :8080")
	http.ListenAndServe(":8080", nil)

}
