package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	http.HandleFunc("/", help)
	http.HandleFunc("/v1/exec", v1Exec)
	fmt.Println("listen on :8080")
	http.ListenAndServe(":8080", nil)
}

func help(w http.ResponseWriter, r *http.Request) {
	helpString := `
curl http://localhost:8080/
curl http://localhost:8080/ping
curl -X POST http://localhost:8080/v1/exec --form-string cmd="ls -al /tmp"
curl -X POST http://localhost:8080/v1/exec --form-string cmd="if true; then echo ok; fi"
curl -X POST http://localhost:8080/v1/exec --form-string cmd="whatever cmd you want to run"

`
	fmt.Fprintf(w, helpString)
}

func v1Exec(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("cmd")

	var cmd *exec.Cmd
	_, err := exec.LookPath("bash")
	if err == nil {
		cmd = exec.Command("bash", "-c", text)
	} else {
		cmd = exec.Command("sh", "-c", text)
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	var result []string
	outScanner := bufio.NewScanner(stdout)
	for outScanner.Scan() {
		result = append(result, outScanner.Text())
	}

	errScanner := bufio.NewScanner(stderr)
	for errScanner.Scan() {
		result = append(result, errScanner.Text())
	}

	fmt.Fprintf(w, strings.Join(result, "\n"))
}
