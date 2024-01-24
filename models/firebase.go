package models

// // import (
// // 	"github.com/beego/beego/v2/client/orm"
// // )

// // type (
// // 	FirebaseFile struct {
// // 		Id       int    `json:"id"`
// // 		FileId   string `json:"file_id"`
// // 		FileLink string `json:"file_link"`
// // 	}
// // )

// // func (t *FirebaseFile) TableName() string {
// // 	return "firebase_file"
// // }

// // func FirebaseFiles() orm.QuerySeter {
// // 	return orm.NewOrm().QueryTable(new(FirebaseFile))
// // }

// // func init() {
// // 	orm.RegisterModel(new(FirebaseFile))
// // }

// // type (
// // 	FirebaseFileRtnJSON struct {
// // 		Id       int    `json:"id"`
// // 		FileId   string `json:"file_id"`
// // 		FileLink string `json:"file_link"`
// // 	}
// // )

// package models

// import (
// 	"time"

// 	"github.com/astaxie/beego/orm"
// )

// type FirebaseFile struct {
// 	Id        int       `json:"id" orm:"column(id);auto;pk"`
// 	FileId    string    `json:"file_id" orm:"column(file_id)"`
// 	FileLink  string    `json:"file_link" orm:"column(file_link)"`
// 	CreatedAt time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
// 	UpdatedAt time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
// 	DeletedAt time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
// }

// func (t *FirebaseFile) TableName() string {
// 	return "firebase_file"
// }

// func FirebaseFiles() orm.QuerySeter {
// 	return orm.NewOrm().QueryTable(new(FirebaseFile))
// }

// func init() {
// 	orm.RegisterModel(new(FirebaseFile))
// }

// func GetAllFirebaseFiles() ([]*FirebaseFile, error) {
// 	var files []*FirebaseFile
// 	_, err := FirebaseFiles().All(&files)
// 	return files, err
// }

// func GetFirebaseFileById(id int) (*FirebaseFile, error) {
// 	fileId := &FirebaseFile{Id: id}
// 	err := orm.NewOrm().Read(fileId)
// 	return fileId, err
// }

// func AddFirebaseFile(file *FirebaseFile) (int64, error) {
// 	return orm.NewOrm().Insert(file)
// }

// func UpdateFirebaseFile(file *FirebaseFile) error {
// 	_, err := orm.NewOrm().Update(file)
// 	return err
// }

// func InsertDataFile(fileLink, fileId string) error {
// 	dataFile := &FirebaseFile{
// 		FileLink: fileLink,
// 		FileId:   fileId,
// 	}
// 	_, err := orm.NewOrm().Insert(dataFile)
// 	return err
// }

// type (
// 	FirebaseFileRtnJSON struct {
// 		Id       int    `json:"id"`
// 		FileId   string `json:"file_id"`
// 		FileLink string `json:"file_link"`
// 	}
// )
