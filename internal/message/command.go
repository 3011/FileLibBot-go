package message

import (
	"fmt"
	"strconv"

	"github.com/3011/FileLibBot-go/internal/config"
	"github.com/3011/FileLibBot-go/internal/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func forwardToGroup(file_id string, filetype string) int {
	var msg tgbotapi.Chattable
	switch filetype {
	case "Document":
		msg = tgbotapi.NewDocument(config.Config.ForwardGroupID, tgbotapi.FileID(file_id))

	case "Video":
		msg = tgbotapi.NewVideo(config.Config.ForwardGroupID, tgbotapi.FileID(file_id))

	case "Photo":
		msg = tgbotapi.NewPhoto(config.Config.ForwardGroupID, tgbotapi.FileID(file_id))
	}
	result, _ := bot.Send(msg)
	return result.MessageID
}

func commandMyFile(formID int64, offset int, editMsgID int) {
	var userFiles []sqlite.User_File
	hasNext := sqlite.FindUserFiles(&userFiles, formID, offset)

	var text string
	if len(userFiles) != 0 {
		text += "*My File:*\n"
		for i, v := range userFiles {
			if i != 0 {
				text += "\n\n"
			}
			text += "[" + v.File_name + "](https://t.me/filelibbot?start=startfile" + v.File_unique_id + ")" + "\n" + v.CreatedAt.Format("2006-01-02 15:04") + "  " + formatFileSize(v.File_size)
		}
	} else {
		text += "Nothing...\nTry sending files to me."
	}

	if editMsgID == 0 {
		msg := tgbotapi.NewMessage(formID, text)
		msg.ParseMode = "Markdown"

		if offset > 0 || hasNext {
			var botton []tgbotapi.InlineKeyboardButton
			if offset > 0 {
				botton = append(botton, tgbotapi.NewInlineKeyboardButtonData("Last", "page last "+strconv.Itoa(offset-1)))
			}
			if hasNext {
				botton = append(botton, tgbotapi.NewInlineKeyboardButtonData("Next", "page next "+strconv.Itoa(offset+1)))
			}
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(botton)
		}

		bot.Send(msg)

	} else {

		var botton []tgbotapi.InlineKeyboardButton
		if offset > 0 || hasNext {

			if offset > 0 {
				botton = append(botton, tgbotapi.NewInlineKeyboardButtonData("Last", "page last "+strconv.Itoa(offset-1)))
			}
			if hasNext {
				botton = append(botton, tgbotapi.NewInlineKeyboardButtonData("Next", "page next "+strconv.Itoa(offset+1)))
			}

		}

		msg := tgbotapi.NewEditMessageTextAndMarkup(formID, editMsgID, text, tgbotapi.NewInlineKeyboardMarkup(botton))
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	}
}

func newDocumentFile(message *tgbotapi.Message) {

	document := message.Document
	fileName := document.FileName
	// 可能 document.FileName 为空
	if fileName == "" {
		fileName = "Document_" + document.FileUniqueID
	}

	msg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
	msg.Caption = "[" + fileName + "](https://t.me/filelibbot?start=file" + message.Document.FileUniqueID + ")" + "\n" + formatFileSize(message.Document.FileSize)
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	file := sqlite.File{File_unique_id: document.FileUniqueID}
	// 判断文件是否存在数据库，没就加上去
	if !sqlite.FindFile(&file) {

		file.File_type = "Document"
		// 转发保存取msgid
		file.Forward_id = forwardToGroup(document.FileID, file.File_type)
		file.File_name = fileName
		file.File_size = document.FileSize
		sqlite.CreateFile(&file)
	}

	userFile := sqlite.User_File{User_id: message.From.ID, File_unique_id: document.FileUniqueID, File_name: fileName, File_size: document.FileSize}
	sqlite.FindUserFile(&userFile)
}

func newVideoFile(message *tgbotapi.Message) {

	video := message.Video
	fileName := video.FileName
	if fileName == "" {
		fileName = "Video_" + video.FileUniqueID
	}

	msg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
	msg.Caption = "[" + fileName + "](https://t.me/filelibbot?start=file" + message.Video.FileUniqueID + ")" + "\n" + formatFileSize(message.Video.FileSize)
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	file := sqlite.File{File_unique_id: video.FileUniqueID}
	if !sqlite.FindFile(&file) {
		// forward
		file.File_type = "Video"
		file.Forward_id = forwardToGroup(video.FileID, file.File_type)
		file.File_name = fileName
		file.File_size = video.FileSize

		sqlite.CreateFile(&file)
	}

	userFile := sqlite.User_File{User_id: message.From.ID, File_unique_id: video.FileUniqueID, File_name: fileName, File_size: video.FileSize}
	sqlite.FindUserFile(&userFile)
}

func newPhotoFile(message *tgbotapi.Message) {
	photo := message.Photo[len(message.Photo)-1]

	fileName := "Photo_" + photo.FileUniqueID
	msg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
	msg.Caption = "[" + fileName + "](https://t.me/filelibbot?start=file" + photo.FileUniqueID + ")" + "\n" + formatFileSize(photo.FileSize)
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	file := sqlite.File{File_unique_id: photo.FileUniqueID}
	if !sqlite.FindFile(&file) {
		// forward
		file.File_type = "Photo"
		file.Forward_id = forwardToGroup(photo.FileID, file.File_type)
		file.File_name = fileName
		file.File_size = photo.FileSize

		sqlite.CreateFile(&file)
	}

	userFile := sqlite.User_File{User_id: message.From.ID, File_unique_id: photo.FileUniqueID, File_name: fileName, File_size: photo.FileSize}
	sqlite.FindUserFile(&userFile)
}

func sendFile(fromId int64, fileUniqueId string) {
	file := sqlite.File{File_unique_id: fileUniqueId}
	if sqlite.FindFile(&file) {
		msg := tgbotapi.NewCopyMessage(fromId, config.Config.ForwardGroupID, file.Forward_id)

		msg.Caption = "[" + file.File_name + "](https://t.me/filelibbot?start=startfile" + file.File_unique_id + ")" + "\n" + formatFileSize(file.File_size)
		msg.ParseMode = "Markdown"

		botton1 := tgbotapi.NewInlineKeyboardButtonData("Delete", "del "+fileUniqueId)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{botton1})
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(fromId, "File not found!")
		bot.Send(msg)
	}
	userFile := sqlite.User_File{User_id: fromId, File_unique_id: file.File_unique_id, File_name: file.File_name, File_size: file.File_size}
	sqlite.FindUserFile(&userFile)
}

func deleteFile(userId int64, fileUniqueID string) {
	sqlite.DelUserFile(userId, fileUniqueID)
}

func formatFileSize(fileSize int) (size string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
