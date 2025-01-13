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
	//"github.com/AlexandrZlnov/go_final_project/service"
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
	//fmt.Println("Таблица создана")
}

func main() {
	EnvirVar := godotenv.Load()
	if EnvirVar != nil {
		log.Fatal("Файл с переменными окружения не найден")
	}

	Port := os.Getenv("TODO_PORT")
	if Port == "" {
		Port = "8080"
		fmt.Println("Переменная среды PORT не установлена, используем порт по умолчанию: 8080")
	}

	// webDir := os.Getenv("WEB_DIR")
	// http.Handle("/", http.FileServer(http.Dir(webDir)))

	NameDBFile := os.Getenv("TODO_DBFILE")
	//fmt.Println(NameDBFile)
	db, err := sql.Open("sqlite", NameDBFile)
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
	r.Delete(("/api/task"), func(w http.ResponseWriter, r *http.Request) { handlers.DeleteTask(w, r, db) } )
	r.Post("/api/task", func(w http.ResponseWriter, r *http.Request) { handlers.PostAddTask(w, r, db) })
	r.Post("/api/task/done", func(w http.ResponseWriter, r *http.Request) {handlers.PostDoneTask(w, r, db) })

	fmt.Println("Server is running on port:", Port)
	if err := http.ListenAndServe(":"+Port, r); err != nil {
		fmt.Printf("Server startup error, %v", err.Error())
		return
	}

	// ----------------------------------------------------------------------------------------

	// загружаем переменные окружения из файла .env
	// errEnv := godotenv.Load()
	// if errEnv != nil {
	// 	log.Fatal("Ошибка при загрузке .env file")
	// }

	// port := os.Getenv("TODO_PORT")
	// if port == "" {
	// 	port = "8080"
	// 	fmt.Println("Переменная среды PORT не установлена, используем порт по умолчанию: 8080")
	// }

	// // запускаем файловый сервер
	// webDir := os.Getenv("WEB_DIR")
	// http.Handle("/", http.FileServer(http.Dir(webDir)))

	// // проверка наличия файла БД
	// // создание фала БД в случае отсутствия
	// db, err := storage.CheckDBFile()
	// if err !=nil {
	// 	log.Fatalf("Ошибка при подключении к БД: %v", err)
	// }
	// defer db.Close()

	// // запускаем роутер Chi
	// r := chi.NewRouter()
	// r.Get("/api/nextdate", handlers.NextDate)
	// r.Get("/index", handlers.GetStatic)
	// r.Post("/api/task", func(w http.ResponseWriter, r *http.Request) {handlers.AddTask(w, r, db)})

	// fmt.Println("Сервер работает на порту:", port)
	// if err := http.ListenAndServe(":"+port, nil); err != nil {
	// 	fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
	// 	return
	// }

}
