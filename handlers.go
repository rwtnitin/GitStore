package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

type Repository struct {
	Name string `json:"name"`
}

func CreateRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	// Read request body which contains the repository name
	var repository Repository
	if err := json.NewDecoder(r.Body).Decode(&repository); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(repository.Name) == "" {
		http.Error(w, "Repository name cannot be empty", http.StatusBadRequest)
		return
	}

	// Create a repository with the name
	var command *exec.Cmd = exec.Command("git", "init", "--bare", fmt.Sprintf("%s.git", repository.Name))
	command.Dir = "data"

	if !DirExists(command.Dir) {
		if err := os.MkdirAll(command.Dir, os.ModePerm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Repository already exists", http.StatusBadRequest)
		return
	}

	if err := command.Run(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the repository name and the repository ID
	w.Write([]byte(fmt.Sprintf("Repository with name '%s' created", repository.Name)))
}
