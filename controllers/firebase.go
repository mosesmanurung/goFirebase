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
	"encoding/json"
	"fmt"
	"io"

	// "json"
	"log"
	"mime"
	"net/http"

	// "os"
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
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	c.CustomAbort(http.StatusInternalServerError, "Error getting current working directory")
	// 	return
	// }
	// filePath := map[string]interface{}{
	// 	"type":                        "service_account",
	// 	"project_id":                  "fir-file-6a929",
	// 	"private_key_id":              "54c1e392f8e7d9c1da7039322291e6855c6a7c82",
	// 	"private_key":                 "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDNTyGFg7wXdNu7\n7feI6zEh92h9j7Y2Vpb08bYHvrfRXJkp7ITjaZZGHxfEO+cvg7XWc8UxkCPbod5d\nevA6q5n/ekbCB4rWTYYBFU2DRkHJMkeYN2LuQwzr4W9/z539Mm1zV387WwaJNwJ7\n1cRSOlu6vo2HJQGhc0vuFKRsBYrQwOJhUlQXwjryUbJpiGeGuTwQxpcXuI0nYD8+\ny+p4ks6tfLSBa3HrxwYwdN3Edcvhwlnv7zrxYM/9qK06461csu0g6D7MdsdSp38u\nLHffQ6BDJiECcYd/zPP3aype+DdnD/S1KyuDxVI9gELP6yO7vaYg+ISFRIVd+x0D\nNiZcJ2RnAgMBAAECggEADX22K7JY4Ue2Crbb7bzauSsrTBjt9csh102s0vx20sSd\ncfJPVyxAijIP2z6+ddJXWBS6cAPTP2L3HDhwYcKV94I+9RAO0P8+H4MZWVd8Ci8K\nIlf9Yb+5MSTasVDgxlsScyJcQ3e7Sbf6K04EPQ0FqxNLdIZ7gXL0mdv0Y/7HPOlZ\ntYZdmlQPEix0PqI2Ecy/vmYvPMAC+kfvlkY5AIUZX/DBX6bxVfC6UZi0OLMgxXfI\n5VHZbYBuBUmlf1gjK533E4OO32RcmJScZikC+oz8dWonofaENlJmVJLZQYzKiqW6\ncHhGFLEkqT7FI5FjuedBrjjR1WYI46/GnR6Ev+CdHQKBgQDmUBNKpF8DqPZFzwhO\nDkSN/An9ALpAa8jFrzF/W/6G3mzXCVcb7PejOYoxVPwff3MrSpquUGvJEF/Ny6L/\nhx1ZOzWxNwwdGi9RhQ2EZvxeL+Vg6dcOvsotqVZojNwXgCniEGIMtay4bw822L3A\nLHRGem382VXvI+qTQQRw6c+grQKBgQDkNSWHBE4PPIuUXfVQISvqrSh08Mj+ZYU7\nqosJEobpYMihnQkSYjLLAPYoQ+GdXlK6L9KWecqU/NHRMrMrVl0nJNHyYCtyIzw6\nsVRiRYtK2Wf7NQn66x745PEhRAiPz5SJ8CY1TppZKqUZBfo4DuqGLBWZKkpLGtic\nuk9R+Ln34wKBgBDXGhIDIs9ps1g3YywR7wFSxIgzXWsIdo292aiuWVYTPXIbxLrO\nAO12b5xb0nObJhisQ9MrHjZ9dAPgN/LnNkYoBi0HEWOvXnZffDWKMjnQ1rzXXFo/\nqRjdoOvUIOO3A1j8Sa8UOaTiugIQpw8+MOJWYWRzn8z0m1pZDrIS5pOhAoGAUmhm\nvtTtI09n0BIF7gOsijgxbdktm8ApVpyFTKdmFIygpsvAZOUVFn2oZm3s4RkYoqd2\nUmR0pUyOsA6w6KttRB9luTLFPZg/vaofoMUgQc00YWCL1BJnwtVZxft9ZAE/0Hfq\nNEHINv7RU7H245tvUThGLGM7JNfy9NXKov1AmqUCgYEAmCEFUFl7RDBogTF2ZrAr\nt8SJl/5XGthm0nlR22+GudkeyddE9/1d6P/PPz9R+3n12ByfrJnekI/mhlLE8P9W\nfxX8mxRPhoYSbigwqV8HtH3G8OOwQ3j4imalJN/rXF8W2jETBPNvNhV+fS/6Povp\nsRiRJk+g9MfxLnSJKojotj4=\n-----END PRIVATE KEY-----\n",
	// 	"client_email":                "firebase-adminsdk-qnpgx@fir-file-6a929.iam.gserviceaccount.com",
	// 	"client_id":                   "100789701306693502943",
	// 	"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
	// 	"token_uri":                   "https://oauth2.googleapis.com/token",
	// 	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	// 	"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-qnpgx%40fir-file-6a929.iam.gserviceaccount.com",
	// 	"universe_domain":             "googleapis.com",
	// }

	// log.Printf("File path: %s, %s", filePath, currentDir)
	filePath, _ := web.AppConfig.String("firebase-storage::firebase_cred")
	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")

	file, header, err := c.GetFile("file")
	fileID := token
	if err != nil {
		log.Printf("File upload error: %v", err)
		c.CustomAbort(http.StatusBadRequest, "File upload error")
		return
	}
	defer file.Close()

	// opt := option.WithCredentialsFile(filePath)
	// jsonData, err := json.Marshal(filePath)
	jsonData, err := json.Marshal(filePath)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v", err)
		c.CustomAbort(http.StatusInternalServerError, fmt.Sprintf("Error marshaling JSON data: %v", err))
		return
	}
	opt := option.WithCredentialsJSON(jsonData)
	config := &firebase.Config{
		StorageBucket: storageBucket,
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

	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", storageBucket, newObjectName, token)

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
