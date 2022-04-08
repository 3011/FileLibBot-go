package callbackquery

import (
	"fmt"
	"strconv"

	"github.com/3011/FileLibBot-go/internal/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func commandMyFile(formID int64, offset int, editMsgID int) {
	// sqlite.CreateFile()
	var userFiles []sqlite.User_File
	hasNext := sqlite.FindUserFiles(&userFiles, formID, offset)
	fmt.Printf("userFiles: %v\n", userFiles)
	// var files []sqlite.File
	// for _, v := range userFiles {
	// 	var file sqlite.File
	// 	file.ID = v.ID
	// 	sqlite.FindFile(&file)
	// 	files = append(files, file)
	// }

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
		// bottonLast := tgbotapi.NewInlineKeyboardButtonData("Last", "page last "+strconv.Itoa(offset-1))

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

func deleteFile(userId int64, fileUniqueID string) {
	sqlite.DelUserFile(userId, fileUniqueID)
}

func formatFileSize(fileSize int) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
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
