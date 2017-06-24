
package main

import (
  "database/sql"
  "fmt"
  "io"
  "log"
  "mime/multipart"
  "net/http"

  _ "github.com/lib/pq"
)

func ping(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "PONG")
}

func createUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
  stmt, err := db.Prepare("INSERT INTO users(first_name, last_name, email, password, user_group_id) VALUES('Mark', 'Crisp', 'mark@sweetiq.com', 'test', 1)")
  if err != nil {
    log.Fatal(err)
  }
  _, err = stmt.Exec()
  if err != nil {
    log.Fatal(err)
  }

  fmt.Fprintf(w, "PONG")
}

func download(w http.ResponseWriter, r *http.Request, db *sql.DB) {
}

func upload(w http.ResponseWriter, r *http.Request, db *sql.DB) {
  if r.Method == "GET" {
    http.Redirect(w, r, "/upload.html", 301)
    } else {
      r.ParseMultipartForm(32 << 20)
      file, handler, err := r.FormFile("uploadfile")
      if err != nil {
        fmt.Println(err)
        return
      }
      defer file.Close()
      http.Redirect(w, r, "/upload.html", 301)
      insertFile(db, file, handler)
    }
  }

  func insertFile(db *sql.DB, file multipart.File, handler *multipart.FileHeader) {
    length, _ := file.Seek(0,2)
    file.Seek(0,0)
    token := make([]byte, length)
    io.ReadFull(file, token)
    stmt, err := db.Prepare("INSERT INTO files(filename, data) VALUES($1, $2)")
    if err != nil {
      log.Fatal(err)
    }
    _, err = stmt.Exec(handler.Filename, token)
    if err != nil {
      log.Fatal(err)
    }
  }

  func main() {
    http.Handle("/", http.FileServer(http.Dir("public")))
    http.HandleFunc("/ping", ping)

    var conninfo string = "postgres://postgres@127.0.0.1:5432?sslmode=disable"
    db, err := sql.Open("postgres", conninfo)
    if err != nil {
      panic(err)
    }

    http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
      upload(w, r, db)
    })
    http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
      download(w, r, db)
    })
    http.HandleFunc("/sign-up", func(w http.ResponseWriter, r *http.Request) {
      createUserHandler(w, r, db)
    })
    http.ListenAndServe(":8080", nil)
  }
