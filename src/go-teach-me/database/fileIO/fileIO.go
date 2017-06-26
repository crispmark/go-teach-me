package fileIO

import (
	"io"
	"log"
	"mime/multipart"
	"time"

	"go-teach-me/database"

	"github.com/google/uuid"
)

type File struct {
	FileID    string
	Filename  string
	CreatedAt time.Time
}

func GetFile(file_id string) (string, []byte) {
	var filename string
	var data []byte
	database.DBCon.QueryRow("SELECT filename, data FROM files WHERE file_id = $1", file_id).Scan(&filename, &data)
	return filename, data
}

func GetAllFileInfo() (*[]File, error) {
	rows, err := database.DBCon.Query("SELECT file_id, filename, created_at FROM files")
	if err != nil {
		return nil, err
	}
	var files []File

	defer rows.Close()
	for rows.Next() {
		var fileID string
		var filename string
		var createdAt time.Time
		err = rows.Scan(&fileID, &filename, &createdAt)
		if err != nil {
			return &files, err
		}
		files = append(files, File{FileID: fileID, Filename: filename, CreatedAt: createdAt})
	}
	return &files, nil
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
