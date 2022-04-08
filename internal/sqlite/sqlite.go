package sqlite

import (
	"github.com/3011/FileLibBot-go/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	User_id         int64 `gorm:"index"`
	File_count      int
	File_size_count int
}

type File struct {
	gorm.Model
	File_unique_id string
	File_name      string
	File_size      int
	Forward_id     int
	File_type      string
}

type User_File struct {
	gorm.Model
	User_id        int64
	File_unique_id string
	File_name      string
	File_size      int
}

func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(config.Config.DBFileName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&User_File{})
}

func FindUser(user *User) {
	result := db.Limit(1).Find(&user, "user_id = ?", user.User_id)

	if result.RowsAffected == 0 {
		db.Create(user)
	}
}

func FindUserFile(userFile *User_File) {
	result := db.Limit(1).Find(&userFile, map[string]interface{}{"user_id": userFile.User_id, "file_unique_id": userFile.File_unique_id})
	if result.RowsAffected == 0 {
		db.Create(userFile)
	}
}

func FindUserFiles(userFiles *[]User_File, userId int64, offset int) bool {
	db.Order("id desc").Offset(offset*8).Limit(8).Find(&userFiles, "user_id = ?", userId)

	var hasNextUserFile []User_File
	result := db.Order("id desc").Offset((offset+1)*8).Limit(1).Find(&hasNextUserFile, "user_id = ?", userId)

	return result.RowsAffected != 0
}

func FindFile(file *File) bool {
	result := db.Limit(1).Find(&file, "file_unique_id = ?", file.File_unique_id)
	return result.RowsAffected != 0
}

func CreateFile(file *File) {
	db.Create(file)
}

func DelUserFile(userId int64, fileUniqueId string) {
	var userFile User_File
	result := db.Limit(1).Find(&userFile, map[string]interface{}{"user_id": userId, "file_unique_id": fileUniqueId})
	if result.RowsAffected != 0 {
		db.Delete(&userFile)
	}
}
