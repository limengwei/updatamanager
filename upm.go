package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	FileDir     = "http://localhost:1024/static/"
	VersionJson = "version.json"
	UploadDir   = "/upload"
)

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	router.Static("/static", "upload")

	router.GET("/", upload)
	router.GET("/upload", upload)
	router.POST("/upload", upload)
	router.GET("/lastversion", lastversion)

	err := router.Run(":1024")

	if err != nil {
		panic(err)
	}

	fmt.Println("server start")
}

//GET：上传文件页面 POST：上传文件操作
func upload(c *gin.Context) {

	if c.Request.Method == "GET" {
		c.HTML(200, "index.html", nil)
	}

	if c.Request.Method == "POST" {
		file, _, err := c.Request.FormFile("apkfile")
		if err != nil {
			message(c, err)
			return
		}
		fileName := strconv.FormatInt(time.Now().Unix(), 16) + ".apk"
		f, err := os.OpenFile(UploadDir+fileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			message(c, err)
			return
		}
		_, err = io.Copy(f, file)
		if err != nil {
			message(c, err)
			return
		}
		defer file.Close()
		defer f.Close()

		//清空内容并写入
		versionFile, err := os.OpenFile(VersionJson, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

		v := Version{}
		v.Name = c.Request.FormValue("versionname")
		v.Code, _ = strconv.Atoi(c.Request.FormValue("versioncode"))
		v.Summary = c.Request.FormValue("versionsummary")
		v.Url = FileDir + fileName

		b, e := json.Marshal(v)
		if e != nil {

		}

		io.WriteString(versionFile, string(b))

		c.String(200, "upload success")
	}
}

//输出错误信息
func message(c *gin.Context, err error) {
	c.HTML(400, "index.html", gin.H{"err": err.Error()})
}

type Version struct {
	Code    int    `json:"version_code"`
	Name    string `json:"version_name"`
	Summary string `json:"version_summary"`
	Url     string `json:"download_url"`
}

//请求最新版本信息
func lastversion(c *gin.Context) {

	file, err := os.Open("version.json")
	if err != nil {

		c.JSON(200, gin.H{
			"ret":  -1,
			"data": "配置文件读取失败",
		})
		return
	}
	defer file.Close()

	v := &Version{}

	dec := json.NewDecoder(file)

	dec.Decode(v)

	c.JSON(200, gin.H{
		"ret":  0,
		"data": v,
	})
}
