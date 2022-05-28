package main

import (
	"bufio"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type Form struct {
	Text string `form:"cmd"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/exec", handler)
	}

	r.Run("0.0.0.0:8080")
}

func handler(c *gin.Context) {
	var form Form
	c.ShouldBind(&form)

	var cmd *exec.Cmd
	_, err := exec.LookPath("bash")
	if err == nil {
		cmd = exec.Command("bash", "-c", form.Text)
	} else {
		cmd = exec.Command("sh", "-c", form.Text)
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		c.String(http.StatusBadRequest, err.Error())
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

	c.String(http.StatusOK, strings.Join(result, "\n"))
}
