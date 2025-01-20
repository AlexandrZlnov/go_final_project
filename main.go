package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AlexandrZlnov/go_final_project/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func CreateDB(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			title TEXT,
			comment TEXT,
			repeat TEXT
		)`)

	if err != nil {
		fmt.Println("Ошибка при создании таблицы:", err)
		return
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);")
	if err != nil {
		fmt.Println("Ошибка при создании индекса:", err)
		return
	}
}

func main() {
	envirVar := godotenv.Load()
	if envirVar != nil {
		log.Fatal("Файл с переменными окружения не найден")
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "8080"
		fmt.Println("Переменная среды PORT не установлена, используем порт по умолчанию: 8080")
	}

	nameDBFile := os.Getenv("TODO_DBFILE")
	db, err := sql.Open("sqlite", nameDBFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install {
		CreateDB(db)
	} else {
		fmt.Println("База данных уже создана")
	}

	r := chi.NewRouter()

	r.Get("/*", handlers.StartFileServer)
	r.Get("/api/nextdate", handlers.GetNextDate)
	r.Get("/api/tasks", func(w http.ResponseWriter, r *http.Request) { handlers.GetTasks(w, r, db) })
	r.Get("/api/task", func(w http.ResponseWriter, r *http.Request) { handlers.GetEditTask(w, r, db) })
	r.Put("/api/task", func(w http.ResponseWriter, r *http.Request) { handlers.PutEditTask(w, r, db) })
	r.Delete(("/api/task"), func(w http.ResponseWriter, r *http.Request) { handlers.DeleteTask(w, r, db) })
	r.Post("/api/task", func(w http.ResponseWriter, r *http.Request) { handlers.PostAddTask(w, r, db) })
	r.Post("/api/task/done", func(w http.ResponseWriter, r *http.Request) { handlers.PostDoneTask(w, r, db) })

	fmt.Println("Server is running on port:", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		fmt.Printf("Server startup error, %v", err.Error())
		return
	}
}
