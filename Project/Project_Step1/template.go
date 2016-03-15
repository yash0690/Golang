//CSCI 130
//Yasasvi Yeleswarapu
//Project Step 1(Create a web application that serves an HTML template)

package main

import (
	"html/template"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("index.html")
	if err != nil {

		log.Fatalln(err) //timestamp

	}
	err = tpl.Execute(w, nil)
}
func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)

}