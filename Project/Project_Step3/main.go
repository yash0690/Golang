//CSCI 130
//Yasasvi Yeleswarapu
//Project Step 3(Create a template which is a form)

package main
import (
		"net/http"						// refers to web server
		"html/template" 					// refers to templates
		"io/ioutil"						// refers to file reading
		"github.com/nu7hatch/gouuid"				// refers to UUID
		"strings"						// refers to split
 		"fmt" 							// refers to debugging		
)
func main() {
	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		name := req.FormValue("name")
		age := req.FormValue("age")
		cookie, err := req.Cookie("session-fino")
		if err != nil {
			id, _ := uuid.NewV4()
			cookie = &http.Cookie{
				Name:  "session-fino",
				Value: id.String() + "|" + name + age,
				// Secure: true,
				HttpOnly: true,
			}
			http.SetCookie(res, cookie)
		}
		err = tpl.Execute(res, nil)
		if err != nil {
			log.Fatalln(err)
		}
	})
	http.ListenAndServe(":8080", nil)
}