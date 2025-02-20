package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

type Task struct {
	Title string
}

func main() {
	// Define the directory where Things 3 stores main.sqlite
	dir := filepath.Join(os.Getenv("HOME"), "Library", "Group Containers", "JLMPQHK86H.com.culturedcode.ThingsMac")

	// Utilize Find to locate the main.sqlite file
	dbPath, err := findMainSQLite(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// If dbPath is not empty, you found the database
	if dbPath != "" {
		fmt.Println("Found database:", dbPath)

		// Now connect to the SQLite database and fetch task titles
		tasks, err := getTaskTitles(dbPath)
		if err != nil {
			fmt.Println("Error fetching tasks:", err)
			return
		}

		// Print all task titles
		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
		} else {
			fmt.Println("Task Titles:")
			for _, task := range tasks {
				fmt.Println(task.Title)
			}
		}
	} else {
		fmt.Println("main.sqlite not found.")
	}
}

// findMainSQLite searches for 'main.sqlite' in the specified directory
func findMainSQLite(dir string) (string, error) {
	var path string
	err := filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(file) == "main.sqlite" {
			path = file
			return filepath.SkipDir // Stop walking once we found the file
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// getTaskTitles opens the database and fetches task titles
func getTaskTitles(dbPath string) ([]Task, error) {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	fmt.Printf("%#v", db)
	// Query to fetch task titles (adjust the table name and column if needed)
	rows, err := db.Query("SELECT TASK.Title FROM TMTASK as TASK") // Replace with the actual table and column names
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Title); err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return tasks, nil
}
