package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

const dataDir string = "data"

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

	// Validate repository name
	if !IsValidRepositoryName(repository.Name) {
		http.Error(w, "Invalid repository name. Only alphabets, numbers, underscores, and hyphens are allowed.", http.StatusBadRequest)
		return
	}

	// Create a repository with the name
	var command *exec.Cmd = exec.Command("git", "init", "--bare", fmt.Sprintf("%s.git", repository.Name))
	command.Dir = dataDir

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

func InfoRefsHandler(w http.ResponseWriter, r *http.Request) {
	var service string = r.URL.Query().Get("service")
	var user string = r.PathValue("user")
	var repository string = r.PathValue("repository")

	var repositoryNamePattern []string = strings.Split(repository, ".")
	if user == "" || repository == "" || len(repositoryNamePattern) < 2 || repositoryNamePattern[1] != "git" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if service == "git-upload-pack" {
		w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
		w.Write([]byte("001e# service=git-upload-pack\n0000"))
		var command *exec.Cmd = exec.Command("git", "upload-pack", "--stateless-rpc", "--advertise-refs", ".")
		command.Dir = path.Join(dataDir, repository)
		command.Stdout = w
		if err := command.Run(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return

	} else if service == "git-receive-pack" {
		fmt.Println("git-receive-pack")
	} else {
		http.Error(w, "Invalid service", http.StatusBadRequest)
		return
	}
}
