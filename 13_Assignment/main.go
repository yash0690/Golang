package main

import (
	"io"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			file, _, err := req.FormFile("my-file")
			if err != nil {
				http.Error(res, err.Error(), 500)
				return
			}
			defer file.Close()

			io.Copy(res, file)

		}

		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, `
      <form method="POST" enctype="multipart/form-data">
        <input type="file" name="my-file">
        <input type="submit">
      </form>
      `)
	}))
}
