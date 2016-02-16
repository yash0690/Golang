//Create the ""surfer page"" and serve it using Go

package main
import (
      "html/template"
      "net/http"
      "log"
)

func surfer(res http.ResponseWriter, req *http.Request)  {
    tmp, err := template.ParseFiles("index.html")
    if err != nil {
        log.Fatalln(err)
    }
    tmp.Execute(res, nil)
}

func main() {
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
    http.Handle("/pic/", http.StripPrefix("/pic/", http.FileServer(http.Dir("pic"))))
    //serving css and img through two statements

    ////http.Handle("/public/",http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
    //serving css and img from one folder that contains both files.

    http.HandleFunc("/", surfer)
    http.ListenAndServe(":8080", nil)
}
