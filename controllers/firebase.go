// package controllers

// import (
// 	"context"
// 	"log"

// 	firebase "firebase.google.com/go/v4"
// 	"firebase.google.com/go/v4/auth"
// 	"google.golang.org/api/option"
// )

// config := &firebase.Config{
// 	StorageBucket: "gs://fir-file-6a929.appspot.com",
// }
// opt := option.WithCredentialsFile("conf\fir-file-6a929-firebase-adminsdk-qnpgx-54c1e392f8.json")
// app, err := firebase.NewApp(context.Background(), config, opt)
// if err != nil {
// 	log.Fatalln(err)
// }

// client, err := app.Storage(context.Background())
// if err != nil {
// 	log.Fatalln(err)
// }

// bucket, err := client.DefaultBucket()
// if err != nil {
// 	log.Fatalln(err)
// }

// firebasefile_controller.go
package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type FirebaseFileController struct {
	web.Controller
}

func (c *FirebaseFileController) Prepare() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, x-xsrf-token, AxiosHeaders, X-Requested-With, X-CSRF-Token, Accept")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Authorization, Set-Cookie, Cookie")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Controller.Prepare()
}

// filePath, _ := web.AppConfig.String("firebase-storage::firebase_cred")
// storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")

func (c *FirebaseFileController) Post() {
	token := uuid.New().String()
	currentDir, err := os.Getwd()
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, "Error getting current working directory")
		return
	}
	filePath := filepath.Join(currentDir, "conf/fir-file-6a929-firebase-adminsdk-qnpgx-54c1e392f8.json")
	log.Printf("File path: %s", filePath)

	file, header, err := c.GetFile("file")
	fileID := token
	if err != nil {
		log.Printf("File upload error: %v", err)
		c.CustomAbort(http.StatusBadRequest, "File upload error")
		return
	}
	defer file.Close()

	opt := option.WithCredentialsFile(filePath)
	config := &firebase.Config{
		StorageBucket: "fir-file-6a929.appspot.com",
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Printf("Firebase app initialization error: %v", err)
		c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Firebase app initialization error: %v", err))
		return
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Printf("Firebase Storage client initialization error: %v", err)
		c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Firebase Storage client initialization error: %v", err))
		return
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Printf("Firebase default bucket retrieval error: %v", err)
		c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Firebase default bucket retrieval error: %v", err))
		return
	}

	newObjectName := fileID
	newObj := bucket.Object(newObjectName)
	wc := newObj.NewWriter(context.Background())

	if _, err := file.Seek(0, 0); err != nil {
		log.Printf("Error while seeking file content: %v", err)
		c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Error while seeking file content: %v", err))
		return
	}
	if _, err := io.Copy(wc, file); err != nil {
		log.Printf("Error while copying file content to Firebase Storage: %v", err)
		c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Error while copying file content to Firebase Storage: %v", err))
		return
	}
	if err := wc.Close(); err != nil {
		log.Printf("Error while closing Firebase Storage writer: %v", err)
		c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Error while closing Firebase Storage writer: %v", err))
		return
	}

	fileType := mime.TypeByExtension(filepath.Ext(header.Filename))

	fmt.Print(fileType)

	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"firebaseStorageDownloadTokens": token,
			"contentType":                   fileType,
		},
	}

	if _, err := bucket.Object(newObjectName).Update(context.Background(), objectAttrsToUpdate); err != nil {
		log.Printf("Error while updating object metadata: %v", err)
		c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Error while updating object metadata: %v", err))
		return
	}

	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", "fir-file-6a929.appspot.com", newObjectName, token)

	c.Data["json"] = map[string]interface{}{"id": 1, "file_id": fileID, "file_link": downloadURL}
	c.ServeJSON()
}

// func (c *FirebaseFileController) PostMulti() {
// 	c.Prepare()
// 	filePath, _ := web.AppConfig.String("firebase-storage::firebase_cred")
// 	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")
// 	pathName := []string{}
// 	pathFile := []string{}

// 	files, err := c.GetFiles("file")

// 	if err != nil {
// 		c.CustomAbort(http.StatusBadRequest, "File upload error")
// 		return
// 	}

// 	for _, fileHeader := range files {
// 		file, err := fileHeader.Open()
// 		if err != nil {
// 			c.CustomAbort(http.StatusBadRequest, "File upload error")
// 			return
// 		}
// 		defer file.Close()

// 		fileID := uuid.New().String()
// 		opt := option.WithCredentialsFile(filePath)
// 		config := &firebase.Config{
// 			StorageBucket: storageBucket,
// 		}
// 		app, err := firebase.NewApp(context.Background(), config, opt)
// 		if err != nil {
// 			log.Fatalln(err)
// 			c.CustomAbort(http.StatusInternalServerError, "Firebase app initialization error")
// 			return
// 		}

// 		client, err := app.Storage(context.Background())
// 		if err != nil {
// 			log.Fatalln(err)
// 			c.CustomAbort(http.StatusInternalServerError, "Firebase Storage client initialization error")
// 			return
// 		}

// 		bucket, err := client.DefaultBucket()
// 		if err != nil {
// 			log.Fatalln(err)
// 			c.CustomAbort(http.StatusInternalServerError, "Firebase default bucket retrieval error")
// 			return
// 		}

// 		newObjectName := fileID
// 		newObj := bucket.Object(newObjectName)
// 		wc := newObj.NewWriter(context.Background())

// 		if _, err := file.Seek(0, 0); err != nil {
// 			log.Fatalln(err)
// 			c.CustomAbort(http.StatusInternalServerError, "Error while seeking file content")
// 			return
// 		}
// 		if _, err := io.Copy(wc, file); err != nil {
// 			log.Fatalln(err)
// 			c.CustomAbort(http.StatusInternalServerError, "Error while copying file content to Firebase Storage")
// 			return
// 		}
// 		if err := wc.Close(); err != nil {
// 			log.Fatalln(err)
// 			c.CustomAbort(http.StatusInternalServerError, "Error while closing Firebase Storage writer")
// 			return
// 		}

// 		fileType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))

// 		fmt.Print(fileType)

// 		objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
// 			Metadata: map[string]string{
// 				"firebaseStorageDownloadTokens": fileID,
// 				"contentType":                   fileType,
// 			},
// 		}

// 		if _, err := bucket.Object(newObjectName).Update(context.Background(), objectAttrsToUpdate); err != nil {
// 			log.Fatalln(err)
// 			c.CustomAbort(http.StatusInternalServerError, "Error while updating object metadata")
// 			return
// 		}

// 		downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", "fir-file-6a929.appspot.com", newObjectName, fileID)
// 		pathFile = append(pathFile, downloadURL)
// 		pathName = append(pathName, fileID)
// 	}

// 	fmt.Print(pathName)

// 	c.Data["json"] = map[string]interface{}{"file_id": pathFile, "file_link": pathName}
// 	c.ServeJSON()
// }
