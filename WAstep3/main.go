package webAssign1

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"html/template"
	"net/http"
)

var tpl *template.Template
var user_list = [3]string{"user1", "user2", "user3"}

const gcsBucket = "csci130-1265.appspot.com"

type gcsPhotos struct {
	ctx    context.Context
	res    http.ResponseWriter
	bucket *storage.BucketHandle
	client *storage.Client
}

// set everything up.
func init() {
	tpl = template.Must(template.ParseGlob("*.html"))
	resourceHandler("css")
	resourceHandler("img")
	resourceHandler("photos")
	http.HandleFunc("/", mainPage)
}

//present main page.
func mainPage(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	gcs := configureCloud(res, req)
	photos := gcs.retrievePhotos()
	if len(photos) == 0 {
		photos = gcs.uploadPhotos() // upload files from local directory.
	}

	user_imgs := make(map[string][]string) //map for passing each user's images to template
	for _, user := range user_list {
		query := &storage.Query{
			Prefix: user + "/",
			// TODO: Delimeter field not needed
			//Delimiter: "/",
		}

		objs, err := gcs.bucket.List(gcs.ctx, query)
		if err != nil {
			log.Errorf(gcs.ctx, "ERROR in query_delimiter")
			return
		}

		var s []string
		//store the public links of the queried result to a slice of string
		for _, obj := range objs.Results {
			s = append(s, obj.MediaLink)
		}
		user_imgs[user] = s
	}
	tpl.Execute(res, user_imgs)
}

// hangle css / img / photo resource directories to point to required resouce files.
func resourceHandler(rsrcDir string) {
	fs := http.FileServer(http.Dir(rsrcDir))
	fs = http.StripPrefix("/"+rsrcDir, fs)
	http.Handle("/"+rsrcDir+"/", fs)
}

// configure cloud context, bucket.
func configureCloud(res http.ResponseWriter, req *http.Request) (gcs *gcsPhotos) {
	ctx := appengine.NewContext(req)
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "ERROR handler NewClient: ", err)
		return
	}
	defer client.Close()

	gcs = &gcsPhotos{
		ctx:    ctx,
		res:    res,
		client: client,
		bucket: client.Bucket(gcsBucket),
	}
	return
}
