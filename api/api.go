package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ChavezJan/dc-final/controller"
	"github.com/gin-gonic/gin"

	// Efectos

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

var info = gin.H{
	"username": gin.H{"email": "username@gmail.com", "token": ""},
}

var tokens = make(map[string]string)

func Start() {

	r := gin.Default()
	r.Use()

	auth := r.Group("/", gin.BasicAuth(gin.Accounts{"username": "password"}))

	auth.GET("/login", login)
	r.DELETE("/logout", logout)
	r.GET("/status", status)
	r.POST("/workloads", workloads)
	r.GET("/workloads/:id", specificWL)
	r.POST("/images", images)
	r.GET("/images/:imgId", download)
	r.Run()

}

func login(c *gin.Context) {

	userToken := c.MustGet(gin.AuthUserKey).(string)

	print(userToken)

	user := c.MustGet(gin.AuthUserKey).(string)
	token := GenerateSecureToken(1)

	tokens[user] = token

	if _, userOk := info[user]; userOk {
		c.JSON(http.StatusOK, gin.H{"message": "Hi " + user + " welcome to the DPIP System", "token": tokens[user]})
	} else {
		c.AbortWithStatus(401)
	}
}

func logout(c *gin.Context) {

	exist, user, _ := auth(c)

	if exist == true {
		delete(tokens, user)
		c.AbortWithStatus(401)
		c.JSON(http.StatusOK, gin.H{"message": "Bye " + user + ", your token has been revoked"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Invalid Token"})
		c.AbortWithStatus(401)
	}
}

func status(c *gin.Context) {

	exist, user, _ := auth(c)

	if exist == true {
		wk := strings.Split(controller.Active_workloads(), "/")
		var wkClean []string

		for i := range wk {
			if wk[i] != "" {
				wkClean = append(wkClean, wk[i])
			}

		}

		current := time.Now()
		c.JSON(http.StatusOK, gin.H{"message": "Hi " + user + ", the DPIP System is Up and Running", "time": current.Format("2006-01-02 15:04:05"), "active": wkClean})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Invalid Token"})
		c.AbortWithStatus(401)
	}
}

func workloads(c *gin.Context) {

	exist, _, _ := auth(c)

	if exist == true {

		filter := c.PostForm("filter")
		WKname := c.PostForm("WKname")
		token := GenerateSecureToken(1)

		if WKname == "" || filter == "" {
			c.AbortWithStatus(401)
			return
		}

		controller.SaveWorkload(WKname, token, true, filter)
		imgId := controller.GetImgIDs(WKname)

		c.JSON(http.StatusOK, gin.H{
			"workload_id":     token,
			"filter":          filter,
			"workload_name":   WKname,
			"status":          true,
			"running_jobs":    10,
			"filtered_images": imgId})

	} else {

		c.JSON(http.StatusOK, gin.H{"message": "Invalid Token"})
		c.AbortWithStatus(401)

	}

}

// no funciona
func specificWL(c *gin.Context) {

	exist, _, _ := auth(c)
	if exist == true {

		name := c.Param("name")
		test := c.GetHeader("data")
		test2, test3 := c.GetQuery("id")

		c.JSON(http.StatusOK, gin.H{"message": "Hi " + name + " - " + test + test2, "bool": test3})

	} else {

		c.JSON(http.StatusOK, gin.H{"message": "Invalid Token"})
		c.AbortWithStatus(401)

	}

}

func images(c *gin.Context) {
	exist, _, _ := auth(c)

	if exist == true {
		image, err := c.FormFile("file")

		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		size := strconv.Itoa(int(image.Size))
		token := GenerateSecureToken(3)
		imgType := "Original"

		controller.SaveImage(image.Filename, token, imgType)

		//----------------
		//modified(img)
		//----------------

		c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "Filename": image.Filename, "filesize": size + " bytes"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Invalid Token"})
		c.AbortWithStatus(401)
	}

}

// no funciona
func modified(imgName string) {

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imgName))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	inverted := effect.Invert(img)
	resized := transform.Resize(inverted, 800, 800, transform.Linear)
	rotated := transform.Rotate(resized, 45, nil)

	//Create a empty file
	file, err := os.Create("./fileName.png")
	if err != nil {
		fmt.Println("ERROR ERROR")
		return
	}
	defer file.Close()

	jpeg.Encode(file, rotated, nil)

	if err := imgio.Save("output.png", rotated, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}

}

// no funciona
func download(c *gin.Context) {

	exist, user, _ := auth(c)
	if exist == true {

		current := time.Now()
		c.JSON(http.StatusOK, gin.H{"message": "Hi " + user + ", the DPIP System is Up and Running", "time": current.Format("2006-01-02 15:04:05"), "active": exist})

	} else {

		c.JSON(http.StatusOK, gin.H{"message": "Invalid Token"})
		c.AbortWithStatus(401)

	}

}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func auth(c *gin.Context) (bool, string, string) {

	exist := false

	bearer := c.Request.Header["Authorization"]
	bearerToken := bearer[0]
	splitedToken := strings.Split(bearerToken, " ")
	token := string(splitedToken[1])

	userName := ""
	userToken := ""

	for user, tokenList := range tokens {

		if token == tokenList {
			exist = true
			userToken = tokenList
			userName = user
		}

	}

	return exist, userName, userToken
}
