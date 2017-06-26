package fileIO

import (
	"io"
	"log"
	"mime/multipart"

	"go-teach-me/database"

	"github.com/google/uuid"
)

func GetFile(file_id string) (string, []byte) {
	var filename string
	var data []byte
	database.DBCon.QueryRow("SELECT filename, data FROM files WHERE file_id = $1", file_id).Scan(&filename, &data)
	return filename, data
}

func InsertFile(file multipart.File, handler *multipart.FileHeader) {
	length, _ := file.Seek(0, 2)
	file.Seek(0, 0)
	token := make([]byte, length)
	io.ReadFull(file, token)
	stmt, err := database.DBCon.Prepare("INSERT INTO files(file_id, filename, data) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(uuid.New().String(), handler.Filename, token)
	if err != nil {
		log.Fatal(err)
	}
}
