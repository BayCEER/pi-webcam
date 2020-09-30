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

var quite bool

func main() {

	flag.BoolVar(&quite, "quite", false, "Turn off pi-webcam's output and error logs")

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

	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/", gin.BasicAuth(gin.Accounts{
		user: password,
	}))

	api.GET("/still", func(c *gin.Context) {
		opts := getOptions(c, stillCommands, stillFlags)
		cmd := exec.Command("raspistill", opts...)
		img, err := cmd.CombinedOutput()
		if err != nil {
			LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.DataFromReader(http.StatusOK, int64(len(img)), "image/"+c.DefaultQuery("encoding", "jpg"), bytes.NewReader(img), nil)
		}
	})
	api.GET("/yuv", func(c *gin.Context) {
		opts := getOptions(c, yuvCommands, yuvFlags)
		cmd := exec.Command("raspiyuv", opts...)
		img, err := cmd.CombinedOutput()
		if err != nil {
			LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.DataFromReader(http.StatusOK, int64(len(img)), "image/octet-stream", bytes.NewReader(img), nil)
		}
	})

	api.GET("/shutdown", func(c *gin.Context) {
		err := exec.Command("sudo", "shutdown", "-h", "now").Start()
		if err != nil {
			LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.String(http.StatusOK, "Going down ...")
		}
	})

	LogInfof("Starting pi-webcam on port:%v\n", port)

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
	LogInfof("Options: %s\n", options)
	return
}

func LogError(err error) {
	if !quite {
		log.Println(err.Error())
	}
}

func LogInfo(s string) {
	if !quite {
		log.Println(s)
	}
}

func LogInfof(format string, v ...interface{}) {
	if !quite {
		log.Printf(format, v...)
	}
}
