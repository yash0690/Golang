package gcsDelimitQry

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"io"
	"net/http"
//	"strings"
        "github.com/nu7hatch/gouuid"
        "encoding/base64"
)

func init() {
	http.HandleFunc("/", handler)
}

const gcsBucket = "csci130-1265.appspot.com"

var quotes = []string {`“Programmers are in a race with the Universe to create bigger and better idiot-proof programs, while the Universe is trying to create bigger and better idiots.  So far the Universe is winning.”
(Rich Cook)`,
`“The trouble with programmers is that you can never tell what a programmer is doing until it’s too late.” (Seymour Cray) `,
`“PHP is a minor evil perpetrated and created by incompetent amateurs, whereas Perl is a great and insidious evil perpetrated by skilled but perverted professionals.”
(Jon Ribbens)`,
`“If debugging is the process of removing bugs, then programming must be the process of putting them in.”
(Edsger W. Dijkstra)`, }

type gcsDelimit struct {
	ctx    context.Context
	res    http.ResponseWriter
	bucket *storage.BucketHandle
	client *storage.Client
}

func handler(res http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	ctx := appengine.NewContext(req)
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "ERROR handler NewClient: ", err)
		return
	}
	defer client.Close()

	d := &gcsDelimit{
		ctx:    ctx,
		res:    res,
		client: client,
		bucket: client.Bucket(gcsBucket),
	}

	d.createFiles()
	d.listFiles()
	io.WriteString(d.res, "\nResults from list directory **with** delimiter\n")
	d.listDir("", "/", "  ")
 	io.WriteString(d.res, "\nResults from list directory **without** delimiter\n")
 	d.listDir("", "", "  ")

}

func (d *gcsDelimit) listDir(name, delim, indent string) {

	query := &storage.Query{
		Prefix:    name,
		Delimiter: delim,
	}

	for query != nil {
		objs, err := d.bucket.List(d.ctx, query)
		if err != nil {
			log.Errorf(d.ctx, "listBucketDirMode: unable to list bucket %q: %v", gcsBucket, err)
			return
		}
		query = objs.Next

		for _, obj := range objs.Results {
			fmt.Fprintf(d.res, "%v%v\n", indent, obj.Name)
		}

		fmt.Fprintf(d.res, "%v\n", objs.Prefixes)

		for _, pfix := range objs.Prefixes {
			log.Infof(d.ctx, "DIR: %v", pfix)
			d.listDir(pfix, delim, indent+"  ")
		}
	}
}

func (d *gcsDelimit) listFiles() {
	io.WriteString(d.res, "\nRetrieving file names...\n")

	client, err := storage.NewClient(d.ctx)
	if err != nil {
		log.Errorf(d.ctx, "%v", err)
		return
	}
	defer client.Close()

	objs, err := client.Bucket(gcsBucket).List(d.ctx, nil)
	if err != nil {
		log.Errorf(d.ctx, "%v", err)
		return
	}

	for _, obj := range objs.Results {
		io.WriteString(d.res, obj.Name+"\n")
	}
}

func (d *gcsDelimit) createFiles() {
	idx := 0
	io.WriteString(d.res, "\nCreating more files for listbucket...\n")
	for _, n := range []string{"foo1", "foo2", "bar", "bar/1", "bar/2", "boo/", "bar/nonce/1", "bar/nonce/2", "bar/nonce/compadre/1", "bar/nonce/compadre/2"} {
		d.createFile(n, idx)
		idx = (idx + 1) & 3
	}
}

func generateUUID(ctx context.Context, fileName string) string {
    uuid, err := uuid.NewV4()
    if err != nil {
		log.Errorf(ctx, "createUUID: unable to write uuid data to bucket %q, file %q: %v", gcsBucket, fileName, err)
    }
    return uuid.String()
}

func (d *gcsDelimit) createFile(fileName string, qid int) {
	fmt.Fprintf(d.res, "Creating file /%v/%v\n", gcsBucket, fileName)

	wc := d.bucket.Object(fileName).NewWriter(d.ctx)
	wc.ContentType = "text/plain"
        b64 := base64.URLEncoding.EncodeToString([]byte(quotes[qid]))
        
        wc.ACL = []storage.ACLRule{
		{storage.AllUsers, storage.RoleReader},
                }

	if _, err := wc.Write([]byte("Header  ")); err != nil {
		log.Errorf(d.ctx, "createFile: unable to write data to bucket %q, file %q: %v", gcsBucket, fileName, err)
		return
	}
	if _, err := wc.Write([]byte(generateUUID(d.ctx, fileName) + "\n\n")); err != nil {
		log.Errorf(d.ctx, "createFile: unable to write data to bucket %q, file %q: %v", gcsBucket, fileName, err)
		return

	}
	if _, err := wc.Write([]byte("Quote:  " + "\n")); err != nil {
		log.Errorf(d.ctx, "createFile: unable to write data to bucket %q, file %q: %v", gcsBucket, fileName, err)
		return
	}
	if _, err := wc.Write([]byte(quotes[qid] + "\n\n")); err != nil {
		log.Errorf(d.ctx, "createFile: unable to write data to bucket %q, file %q: %v", gcsBucket, fileName, err)
		return
	}

	if _, err := wc.Write([]byte("Quote encoded  " + "\n")); err != nil {
		log.Errorf(d.ctx, "createFile: unable to write data to bucket %q, file %q: %v", gcsBucket, fileName, err)
		return
	}
	if _, err := wc.Write([]byte(b64 + "\n\n")); err != nil {
		log.Errorf(d.ctx, "createFile: unable to write data to bucket %q, file %q: %v", gcsBucket, fileName, err)
		return
	}
	if err := wc.Close(); err != nil {
		log.Errorf(d.ctx, "createFile: unable to close bucket %q, file %q: %v", gcsBucket, fileName, err)
		return
	}
}

