package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"os"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"//json封装解析"
)



func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/generateQRCode/:key/:token", func(c *gin.Context) {
		qr_id := c.Params.ByName("key")
		token := c.Params.ByName("token")
		rs,tf:=getQRCode(qr_id,token)
		if tf{
		   c.JSON(200, gin.H{"qrcode":rs,"error":0})
		}else{
		   c.JSON(200, gin.H{"error":1})	
		}
	})
	r.GET("/deleteQRCode/:key", func(c *gin.Context) {
		qr_id := c.Params.ByName("key")
		tf:=deleteQRCode(qr_id)
		if tf{
		   c.JSON(200, gin.H{"error":0})
		}else{
		   c.JSON(200, gin.H{"error":1})	
		}
	})
	return r
}

func getQRCode(key string,token string)(string,bool){
	tf:=false
	qrurl:= "assets/qrcode/"+key+".jpg"
	var buf bytes.Buffer
	buf.WriteString("https://api.weixin.qq.com/wxa/getwxacodeunlimit")
	buf.WriteString("?access_token="+token)
	urlstr:=buf.String()
	width:=430
	params:= make(map[string]interface{})
	params["scene"] = key
    params["width"] = width
    params["auto_color"] = true
    bytesData, err := json.Marshal(params)
    if err != nil {
        fmt.Println(err.Error() )
        return "",tf
	}
	fmt.Println("json_str:",string(bytesData))
	reader := bytes.NewReader(bytesData)
    request, err := http.NewRequest("POST", urlstr, reader)
	if err != nil {
	   fmt.Println(err.Error())
	   return "",tf
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
    resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
        return "",tf
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        fmt.Println(err.Error())
        return "",tf
	}
	img, _, _ := image.Decode(bytes.NewReader(respBytes))
	
	out, err := os.Create(qrurl)
	if err != nil {
		fmt.Println(err)
		return "",tf
	}
	err = jpeg.Encode(out,img,nil)
	if err != nil {
		fmt.Println(err)
		return "",tf
	}
	return "qrcode/"+key+".jpg",true
}

func deleteQRCode(key string)bool{
	tf:=true
	filePath:= "assets/qrcode/"+key+".jpg"
	err := os.Remove(filePath)
	if err != nil{
		tf=false
	}
	return tf
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	//r.RunTLS("192.168.6.103:8080","214801461110100.pem","214801461110100.key")
	r.Run("127.0.0.1:8080")
}
	
