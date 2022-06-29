package main

import (
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	root := os.Getenv("GIT_ROOT")

	username := os.Getenv("GIT_USERNAME")
	password := os.Getenv("GIT_PASSWORD")

	s, err := New(root, username, password)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Git Server started on :8080")
	log.Fatal(s.ListenAndServe(":8080"))
}

type Server struct {
	root string

	username string
	password string

	pathGit string
}

func New(path, username, password string) (*Server, error) {
	root, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	pathGit, err := exec.LookPath("git")

	if err != nil {
		return nil, err
	}

	s := &Server{
		root: root,

		username: username,
		password: password,

		pathGit: pathGit,
	}

	return s, nil
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()

	if s.username != username || s.password != password {
		w.Header().Set("WWW-Authenticate", "Basic realm=git")
		w.WriteHeader(401)
		return
	}

	if username == "" {
		username = "git"
	}

	const prefixAdminRepo = "/admin/repo/"

	if strings.HasPrefix(r.URL.Path, prefixAdminRepo) {
		repo := strings.TrimPrefix(r.URL.Path, prefixAdminRepo)
		repoPath := path.Join(s.root, repo)

		if r.Method == http.MethodPost {
			if err := exec.Command("git", "init", "--bare", repoPath).Run(); err != nil {
				w.WriteHeader(400)
				return
			}

			w.WriteHeader(201)
			return
		}

		if r.Method == http.MethodDelete {
			if err := os.RemoveAll(repoPath); err != nil {
				w.WriteHeader(400)
				return
			}

			w.WriteHeader(200)
			return
		}
	}

	h := cgi.Handler{
		Dir: s.root,

		Env: []string{
			"REMOTE_USER=" + username,
			"GIT_HTTP_EXPORT_ALL=",
			"GIT_PROJECT_ROOT=" + s.root,
		},

		Path: s.pathGit,
		Args: []string{
			"http-backend",
		},
	}

	if len(r.TransferEncoding) > 0 && r.TransferEncoding[0] == `chunked` {
		r.TransferEncoding = nil
		r.Header.Set(`Transfer-Encoding`, `chunked`)

		r.ContentLength = -1
	}

	h.ServeHTTP(w, r)
}
