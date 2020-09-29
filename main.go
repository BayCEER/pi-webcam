package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

var stillCommands = []string{
	"encoding",
	"exif",
	"quality",

	"awb",
	"awbgains",
	"timeout",
	"exposure",
	"metering",
	"flicker",
	"drc",
	"sharpness",
	"contrast",
	"brightness",
	"saturation",
	"ISO",
	"shutter",
}

var stillFlags = []string{
	"hflip",
	"vflip",
	"raw",
}

var yuvCommands = []string{
	"awb",
	"awbgains",
	"timeout",
	"exposure",
	"metering",
	"flicker",
	"drc",
	"sharpness",
	"contrast",
	"brightness",
	"saturation",
	"ISO",
	"shutter",
}

var yuvFlags = []string{
	"rgb",
	"bgr",
	"luma",
	"hflip",
	"vflip",
}

func main() {

	var user string
	flag.StringVar(&user, "user", "webcam", "User name to access the camera")

	var password string
	flag.StringVar(&password, "password", "webcam", "Password to access the camera")

	var port int
	flag.IntVar(&port, "port", 8080, "Server port")

	var tls bool
	flag.BoolVar(&tls, "tls", false, "Secured by TLS, needs certFile and keyFile options to be set.")

	// HTTPS Settings
	var certFile string
	var keyFile string
	flag.StringVar(&certFile, "certFile", "./cert/certfile.crt", "Path to cert file for TLS")
	flag.StringVar(&keyFile, "keyFile", "./cert/keyfile.key", "Path to private key file for TLS")

	flag.Parse()

	r := gin.Default()

	api := r.Group("/", gin.BasicAuth(gin.Accounts{
		user: password,
	}))

	api.GET("/still", func(c *gin.Context) {
		opts := getOptions(c, stillCommands, stillFlags)
		cmd := exec.Command("raspistill", opts...)
		img, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			log.Println("Image captured")
			c.DataFromReader(http.StatusOK, int64(len(img)), "image/"+c.DefaultQuery("encoding", "jpg"), bytes.NewReader(img), nil)
		}
	})
	api.GET("/yuv", func(c *gin.Context) {
		opts := getOptions(c, yuvCommands, yuvFlags)
		cmd := exec.Command("raspiyuv", opts...)
		img, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			log.Println("Image captured")
			c.DataFromReader(http.StatusOK, int64(len(img)), "image/octet-stream", bytes.NewReader(img), nil)
		}
	})

	log.Printf("Starting pi-webcam on port:%v\n", port)

	if tls {
		r.RunTLS(fmt.Sprintf(":%v", port), certFile, keyFile)
	} else {
		r.Run(fmt.Sprintf(":%v", port))
	}

}

func getOptions(c *gin.Context, commands []string, flags []string) (options []string) {
	options = []string{"-o", "-"}
	for _, command := range commands {
		if value, r := c.GetQuery(command); r {
			options = append(options, "--"+command, value)
		}
	}
	for _, flag := range flags {
		if _, r := c.GetQuery(flag); r {
			options = append(options, "--"+flag)
		}
	}
	log.Printf("Options: %s\n", options)
	return
}
