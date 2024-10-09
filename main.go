package main

import (
        "encoding/json"
        "net/http"
        "strconv"

        "database/sql"
        _ "github.com/mattn/go-sqlite3"
)

type Task struct {
        ID      int    `json:"id"`
        Title   string `json:"title"`
        Details string `json:"details"`
        Done    bool   `json:"done"`
}

var db *sql.DB

func init() {
        var err error
        db, err = sql.Open("sqlite3", "tasks.db")
        if err != nil {
            panic(err)
        }
        _, err = db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, details TEXT, done BOOLEAN)")
        if err != nil {
            panic(err)
        }
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query("SELECT id, title, details, done FROM tasks")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var tasks []Task
        for rows.Next() {
            var task Task
            err = rows.Scan(&task.ID, &task.Title, &task.Details, &task.Done)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            tasks = append(tasks, task)
        }

        json.NewEncoder(w).Encode(tasks)
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.URL.Query().Get("id"))
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }

        var task Task
        err = db.QueryRow("SELECT id, title, details, done FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Title, &task.Details, &task.Done)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.URL.Query().Get("id"))
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }

        var task Task
        err = json.NewDecoder(r.Body).Decode(&task)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        _, err = db.Exec("UPDATE tasks SET title = ?, details = ?, done = ? WHERE id = ?", task.Title, task.Details, task.Done, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.URL.Query().Get("id"))
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }

        _, err = db.Exec("DELETE FROM tasks WHERE id = ?", id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
}

func main() {
        http.HandleFunc("/tasks", getAllTasks)
        http.HandleFunc("/tasks/get", getTaskById)
        http.HandleFunc("/tasks/update", updateTask)
        http.HandleFunc("/tasks/delete", deleteTask)

        http.ListenAndServe(":8080", nil)
}