package myapp

import (
	"crypto/sha1"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"google.golang.org/cloud/storage"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine"
	"golang.org/x/net/context"
)


func uploadPhoto(src multipart.File, hdr *multipart.FileHeader, c *http.Cookie, req *http.Request) *http.Cookie {


	
	defer src.Close()


	ctx := appengine.NewContext(req)


	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:]




	switch ext {
	case "jpg", "jpeg", "png", "txt":
		//log.Infof(ctx, "GOOD FILE EXTENSION: %s", ext)
	default:
		log.Errorf(ctx, "We do not allow files of type %s. We only allow jpg, jpeg, png, txt extensions.", ext)
		return c
	}


	m := Model(c)

	var fName string


	if(ext == "jpeg" || ext == "jpg" || ext == "png"){
		fName = m.Name + "/image/"
	}else{
		fName = m.Name + "/text/"
	}

	
	fName = fName + getSha(src) + "." + ext

	src.Seek(0, 0)



	client, err := storage.NewClient(ctx)
	if err != nil{
		log.Errorf(ctx, "Error in main client err")
		return c
	}
	defer client.Close()


	writer := client.Bucket(gcsBucket).Object(fName).NewWriter(ctx)


	
	writer.ACL = []storage.ACLRule{
		{storage.AllUsers, storage.RoleReader},
	}

	
	if(ext == "jpeg" || ext == "jpg") {
		writer.ContentType = "image/jpeg"
	}else if(ext == "png"){
		writer.ContentType = "image/png"
	}else{
		writer.ContentType = "text/plain"
	}


	io.Copy(writer, src)

	err = writer.Close()
	if(err != nil){
		log.Errorf(ctx, "error uploadPhoto writer close", err)
		return c
	}

	fName, err = getFileLink(ctx, fName)

	return addPhoto(fName, ext, c)
}


func addPhoto(fName string, ext string, c *http.Cookie) *http.Cookie {

	m := Model(c)

	if ext == "jpg" || ext == "png" || ext == "jpeg"{
		m.Pictures = append(m.Pictures, fName)
	}
	//Store the file path in string slice, m.Files
	m.Files = append(m.Files, fName)
	//Get id from old Model and update the cookie with updated model
	xs := strings.Split(c.Value, "|")
	id := xs[0]
	cookie := currentVisitor(m, id)
	return cookie
}

func getSha(src multipart.File) string {
	h := sha1.New()
	io.Copy(h, src)
	return fmt.Sprintf("%x", h.Sum(nil))
}


//returns the download link
func getFileLink(ctx context.Context, name string) (string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	attrs, err := client.Bucket(gcsBucket).Object(name).Attrs(ctx)
	if err != nil {
		return "", err
	}
	return attrs.MediaLink, nil
}
