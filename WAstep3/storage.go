package webAssign1

import (
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
	"math/rand"
	"time"
)

// load photos from the cloud to memory.
func (gcs *gcsPhotos) retrievePhotos() (photos []string) {
	files, err := gcs.bucket.List(gcs.ctx, nil)
	if err != nil {
		log.Errorf(gcs.ctx, "listBucketDirMode: unable to list bucket %q: %v", gcsBucket, err)
		return
	}
	for _, name := range files.Results {
		photos = append(photos, name.Name)
	}
	return
}

// write image file to GCS.
func (gcs *gcsPhotos) write(filename string, data []byte) {
	writer := gcs.bucket.Object(filename).NewWriter(gcs.ctx)
	writer.ACL = []storage.ACLRule{
		{storage.AllUsers, storage.RoleReader},
	}
	writer.ContentType = "image/jpg"
	defer writer.Close()
	writer.Write(data)
}

// up load photos from local location to the cloud.
func (gcs *gcsPhotos) uploadPhotos() []string {
	ext := ".jpg"
	subDir := "photos/"
	testPhotos := []string{"0", "2", "3", "4", "5", "6", "7", "8", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24"}

	for _, name := range testPhotos {
		ffile := subDir + name + ext
		srcFile, err := gcs.read(ffile)
		if err != nil {
			log.Errorf(gcs.ctx, "Open / read file error %v", err)
			return nil
		}
		fName := gcs.randomUser() + "photo" + name + ext
		gcs.write(fName, srcFile)
		log.Infof(gcs.ctx, "In File: %s, Out File: %s\n", ffile, fName)
	}
	return gcs.retrievePhotos()
}

func (gcs *gcsPhotos) randomUser() (user string) {
	rand.Seed(time.Now().Unix())
	x := rand.Intn(3)
	user = user_list[x] + "/"
	return
}
