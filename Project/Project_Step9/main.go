package main

import (
	"net/http"
	
	"encoding/json"
	"encoding/base64"
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"strings"
)

type User struct {
	Name string
	Age string
	Log string
}

func main() {
	http.HandleFunc("/", foo)
	http.HandleFunc("/login",textFile)
	http.ListenAndServe(":8080", nil)
}

func getCode(data string) string {
	h := hmac.New(sha256.New, []byte("raka02asd2j01djmk"))
	io.WriteString(h,data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func foo(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("logged-in")
	//no cookie
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "logged-in",
			Value: "0",
			//Secure: true,
			HttpOnly: true,
		}
	}


	if req.Method == "POST" {

		password := req.FormValue("password")

		if password == "secret" {
			person := User{
				Name: req.FormValue("userName"),
				Age: req.FormValue("age"),
				Log: "1",
			}
			b, _ := json.Marshal(person)
			encode := base64.StdEncoding.EncodeToString(b)
			id,_ := uuid.NewV4()

			cookie = &http.Cookie{
				Name:  "logged-in",
				Value: "1" + "|" + id.String() + "|" + getCode(id.String()) + "|" + encode + "|" + getCode(encode),
				//Secure: true,
				HttpOnly: true,
			}
		}
	}

	if req.URL.Path == "/logout" {
		cookie = &http.Cookie{
			Name:   "logged-in",
			Value:  "0",
			MaxAge: -1,
			
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
		http.Redirect(res, req, "/", 302)
	}
	http.SetCookie(res, cookie)
	cv := strings.Split(cookie.Value, "|")

	var html string


	if cv[0] == "0" {
		html = `
				<!DOCTYPE html>
				<html lang="en">
				<head>
						<meta charset="UTF-8">
						<title></title>
				</head>
				<body>
				<h1>LOG IN</h1>
				<form method= "Post">
						<h4>User Name:</h4><input type="text" name="userName"><br>
						<h4>Age:<h4><input type="text" name="age"><br>
						<h4>Password:</h4><input type="text" name="password"><br>
						<input type="submit">
				</form>
				</body>
				</html>`
	}
	//logged in
	if cv[0] == "1" {
		http.Redirect(res, req, "/login", 302)


	}
	io.WriteString(res, html)
}



func textFile(res http.ResponseWriter,req *http.Request){
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprint(res,`
		<h1><a href= "/logout">LOG OUT</a></h1>
		<form method = "POST"  enctype="multipart/form-data">
		<input type="file" name="name"><br><br>
		<input type="submit">
		</form>`)
	if req.Method == "POST"{
		key := "name"
		_, header, err := req.FormFile(key)
		if err != nil{
			fmt.Println(err)
		}
		inFile, err := header.Open()
		if err != nil{
			fmt.Println(err)
		}
		io.Copy(res,inFile)
	}
}
