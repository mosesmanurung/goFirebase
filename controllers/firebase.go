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
	"path/filepath"
	"strings"

	// "strings"

	// "path/filepath"

	// "strconv"

	"EasyGo/models"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type FirebaseFileController struct {
	web.Controller
}

func (c *FirebaseFileController) Post() {
	token := uuid.New().String()
	filePath, _ := web.AppConfig.String("firebase-storage::firebase_cred")
	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")

	file, header, err := c.GetFile("file")
	fileID := token
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "File upload error")
		return
	}
	defer file.Close()

	opt := option.WithCredentialsFile(filePath)
	config := &firebase.Config{
		StorageBucket: storageBucket,
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Firebase app initialization error")
		return
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Firebase Storage client initialization error")
		return
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Firebase default bucket retrieval error")
		return
	}

	newObjectName := fileID
	newObj := bucket.Object(newObjectName)
	wc := newObj.NewWriter(context.Background())

	if _, err := file.Seek(0, 0); err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while seeking file content")
		return
	}
	if _, err := io.Copy(wc, file); err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while copying file content to Firebase Storage")
		return
	}
	if err := wc.Close(); err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while closing Firebase Storage writer")
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
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while updating object metadata")
		return
	}

	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", "fir-file-6a929.appspot.com", newObjectName, token)

	firebaseFile := models.FirebaseFile{
		FileId:   fileID,
		FileLink: downloadURL,
	}

	fmt.Println(firebaseFile)

	fileId, err := models.AddFirebaseFile(&firebaseFile)
	if err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while adding FirebaseFile to the database")
		return
	}

	c.Data["json"] = map[string]interface{}{"id": fileId, "file_id": fileID, "file_link": downloadURL}
	c.ServeJSON()
}

func (c *FirebaseFileController) PostMulti() {
	// token := uuid.New().String()
	filePath, _ := web.AppConfig.String("firebase-storage::firebase_cred")
	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")
	pathName := []string{}
	pathFile := []string{}

	// file, header, err := c.GetFile("file")
	files, err := c.GetFiles("file")

	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "File upload error")
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.CustomAbort(http.StatusBadRequest, "File upload error")
			return
		}
		defer file.Close()

		fileID := uuid.New().String()
		opt := option.WithCredentialsFile(filePath)
		config := &firebase.Config{
			StorageBucket: storageBucket,
		}
		app, err := firebase.NewApp(context.Background(), config, opt)
		if err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Firebase app initialization error")
			return
		}

		client, err := app.Storage(context.Background())
		if err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Firebase Storage client initialization error")
			return
		}

		bucket, err := client.DefaultBucket()
		if err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Firebase default bucket retrieval error")
			return
		}

		newObjectName := fileID
		newObj := bucket.Object(newObjectName)
		wc := newObj.NewWriter(context.Background())

		if _, err := file.Seek(0, 0); err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Error while seeking file content")
			return
		}
		if _, err := io.Copy(wc, file); err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Error while copying file content to Firebase Storage")
			return
		}
		if err := wc.Close(); err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Error while closing Firebase Storage writer")
			return
		}

		fileType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))

		fmt.Print(fileType)

		objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
			Metadata: map[string]string{
				"firebaseStorageDownloadTokens": fileID,
				"contentType":                   fileType,
			},
		}

		if _, err := bucket.Object(newObjectName).Update(context.Background(), objectAttrsToUpdate); err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Error while updating object metadata")
			return
		}

		downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", "fir-file-6a929.appspot.com", newObjectName, fileID)
		pathFile = append(pathFile, downloadURL)
		pathName = append(pathName, fileID)
	}

	fmt.Print(pathName)

	firebaseFile := models.FirebaseFile{
		FileId:   strings.Join(pathName, ","),
		FileLink: strings.Join(pathFile, ","),
	}

	fmt.Println(firebaseFile)

	fileId, err := models.AddFirebaseFile(&firebaseFile)
	if err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while adding FirebaseFile to the database")
		return
	}

	c.Data["json"] = map[string]interface{}{"id": fileId, "file_id": pathFile, "file_link": pathName}
	c.ServeJSON()
}

func (c *FirebaseFileController) PostLink() {
	fileId := c.GetString("file_id")
	fileLink := c.GetString("file_link")

	err := models.InsertDataFile(fileLink, fileId)
	if err != nil {
		fmt.Print(err)
		c.Data["json"] = map[string]interface{}{"error": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"status": "success"}
	}
	c.ServeJSON()
}

func (c *FirebaseFileController) GetAll() {
	files, err := models.GetAllFirebaseFiles()
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, "Error retrieving FirebaseFiles")
		return
	}

	c.Data["json"] = files
	c.ServeJSON()
}

func (c *FirebaseFileController) GetOne() {
	id, _ := c.GetInt(":id")

	file, err := models.GetFirebaseFileById(id)
	if err != nil {
		c.CustomAbort(http.StatusNotFound, "FirebaseFile not found")
		return
	}

	c.Data["json"] = file
	c.ServeJSON()
}

// func (c *FirebaseFileController) Put() {
// 	token := uuid.New().String()
// 	fileID := token
// 	id, _ := c.GetInt(":id")
// 	file, header, err := c.GetFile("file")

// 	if err != nil {
// 		c.CustomAbort(http.StatusBadRequest, "File upload error")
// 		return
// 	}
// 	defer file.Close()

// 	filePath, _ := web.AppConfig.String("firebase-storage::firebase_cred")
// 	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")

// 	opt := option.WithCredentialsFile(filePath)
// 	config := &firebase.Config{
// 		StorageBucket: storageBucket,
// 	}
// 	app, err := firebase.NewApp(context.Background(), config, opt)
// 	if err != nil {
// 		log.Fatalln(err)
// 		c.CustomAbort(http.StatusInternalServerError, "Firebase app initialization error")
// 		return
// 	}

// 	client, err := app.Storage(context.Background())
// 	if err != nil {
// 		log.Fatalln(err)
// 		c.CustomAbort(http.StatusInternalServerError, "Firebase Storage client initialization error")
// 		return
// 	}

// 	bucket, err := client.DefaultBucket()
// 	if err != nil {
// 		log.Fatalln(err)
// 		c.CustomAbort(http.StatusInternalServerError, "Firebase default bucket retrieval error")
// 		return
// 	}

// 	newObjectName := fileID
// 	newObj := bucket.Object(newObjectName)
// 	wc := newObj.NewWriter(context.Background())

// 	if _, err := io.Copy(wc, file); err != nil {
// 		log.Fatalln(err)
// 		c.CustomAbort(http.StatusInternalServerError, "Error while copying file content to Firebase Storage")
// 		return
// 	}

// 	if err := wc.Close(); err != nil {
// 		log.Fatalln(err)
// 		c.CustomAbort(http.StatusInternalServerError, "Error while closing Firebase Storage writer")
// 		return
// 	}

// 	fileType := mime.TypeByExtension(filepath.Ext(header.Filename))

// 	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
// 		Metadata: map[string]string{
// 			"firebaseStorageDownloadTokens": token,
// 			"contentType":                   fileType,
// 		},
// 	}

// 	if _, err := bucket.Object(newObjectName).Update(context.Background(), objectAttrsToUpdate); err != nil {
// 		log.Fatalln(err)
// 		c.CustomAbort(http.StatusInternalServerError, "Error while updating object metadata")
// 		return
// 	}

// 	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", storageBucket, newObjectName, token)

// 	firebaseFile := models.FirebaseFile{
// 		Id:       id,
// 		FileId:   fileID,
// 		FileLink: downloadURL,
// 	}

// 	if err := models.UpdateFirebaseFile(&firebaseFile); err != nil {
// 		log.Fatalln(err)
// 		c.CustomAbort(http.StatusInternalServerError, "Error while updating FirebaseFile in the database")
// 		return
// 	}

// 	c.Data["json"] = map[string]interface{}{"id": id, "file_id": fileID, "file_link": downloadURL}
// 	c.ServeJSON()
// }

func (c *FirebaseFileController) Put() {
	id, _ := c.GetInt(":id")

	existingFile, err := models.GetFirebaseFileById(id)
	if err != nil {
		c.CustomAbort(http.StatusNotFound, "File not found")
		return
	}

	err = deleteFileFromStorage(existingFile.FileId)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, "Error deleting file from Firebase Storage")
		return
	}

	token := uuid.New().String()
	fileID := token

	file, header, err := c.GetFile("file")
	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "File upload error")
		return
	}
	defer file.Close()

	filePath, _ := web.AppConfig.String("firebase-storage::firebase_cred")
	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")

	opt := option.WithCredentialsFile(filePath)
	config := &firebase.Config{
		StorageBucket: storageBucket,
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Firebase app initialization error")
		return
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Firebase Storage client initialization error")
		return
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Firebase default bucket retrieval error")
		return
	}

	newObjectName := fileID
	newObj := bucket.Object(newObjectName)
	wc := newObj.NewWriter(context.Background())

	if _, err := io.Copy(wc, file); err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while copying file content to Firebase Storage")
		return
	}

	if err := wc.Close(); err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while closing Firebase Storage writer")
		return
	}

	fileType := mime.TypeByExtension(filepath.Ext(header.Filename))

	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"firebaseStorageDownloadTokens": token,
			"contentType":                   fileType,
		},
	}

	if _, err := bucket.Object(newObjectName).Update(context.Background(), objectAttrsToUpdate); err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while updating object metadata")
		return
	}

	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", storageBucket, newObjectName, token)

	firebaseFile := models.FirebaseFile{
		Id:       id,
		FileId:   fileID,
		FileLink: downloadURL,
	}

	if err := models.UpdateFirebaseFile(&firebaseFile); err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while updating FirebaseFile in the database")
		return
	}

	c.Data["json"] = map[string]interface{}{"id": id, "file_id": fileID, "file_link": downloadURL}
	c.ServeJSON()
}

func (c *FirebaseFileController) PutMulti() {
	id, _ := c.GetInt(":id")
	filePath, _ := web.AppConfig.String("firebase-storage::firebase_cred")
	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")
	pathName := []string{}
	pathFile := []string{}

	existingFile, err := models.GetFirebaseFileById(id)

	if err != nil {
		c.CustomAbort(http.StatusNotFound, "File not found")
		return
	}

	filesd := strings.Split(existingFile.FileId, ",")

	for _, temp := range filesd {
		err = deleteFileFromStorage(temp)
		if err != nil {
			c.CustomAbort(http.StatusInternalServerError, "Error deleting file from Firebase Storage")
			return
		}
	}

	files, err := c.GetFiles("file")

	if err != nil {
		c.CustomAbort(http.StatusBadRequest, "File upload error")
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.CustomAbort(http.StatusBadRequest, "File upload error")
			return
		}
		defer file.Close()

		fileID := uuid.New().String()
		opt := option.WithCredentialsFile(filePath)
		config := &firebase.Config{
			StorageBucket: storageBucket,
		}
		app, err := firebase.NewApp(context.Background(), config, opt)
		if err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Firebase app initialization error")
			return
		}

		client, err := app.Storage(context.Background())
		if err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Firebase Storage client initialization error")
			return
		}

		bucket, err := client.DefaultBucket()
		if err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Firebase default bucket retrieval error")
			return
		}

		newObjectName := fileID
		newObj := bucket.Object(newObjectName)
		wc := newObj.NewWriter(context.Background())

		if _, err := file.Seek(0, 0); err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Error while seeking file content")
			return
		}
		if _, err := io.Copy(wc, file); err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Error while copying file content to Firebase Storage")
			return
		}
		if err := wc.Close(); err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Error while closing Firebase Storage writer")
			return
		}

		fileType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))

		fmt.Print(fileType)

		objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
			Metadata: map[string]string{
				"firebaseStorageDownloadTokens": fileID,
				"contentType":                   fileType,
			},
		}

		if _, err := bucket.Object(newObjectName).Update(context.Background(), objectAttrsToUpdate); err != nil {
			log.Fatalln(err)
			c.CustomAbort(http.StatusInternalServerError, "Error while updating object metadata")
			return
		}

		downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", "fir-file-6a929.appspot.com", newObjectName, fileID)
		pathFile = append(pathFile, downloadURL)
		pathName = append(pathName, fileID)
	}

	fmt.Print(pathName)

	firebaseFile := models.FirebaseFile{
		FileId:   strings.Join(pathName, ","),
		FileLink: strings.Join(pathFile, ","),
	}

	fmt.Println(firebaseFile)

	fileId, err := models.AddFirebaseFile(&firebaseFile)
	if err != nil {
		log.Fatalln(err)
		c.CustomAbort(http.StatusInternalServerError, "Error while adding FirebaseFile to the database")
		return
	}

	c.Data["json"] = map[string]interface{}{"id": fileId, "file_id": pathFile, "file_link": pathName}
	c.ServeJSON()
}

func deleteFileFromStorage(fileID string) error {
	filePath, _ := web.AppConfig.String("firebase-storage::firebase_cred")
	storageBucket, _ := web.AppConfig.String("firebase-storage::bucket_link")

	opt := option.WithCredentialsFile(filePath)
	config := &firebase.Config{
		StorageBucket: storageBucket,
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
		return err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	err = bucket.Object(fileID).Delete(context.Background())
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
