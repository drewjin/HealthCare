package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ImageOcr(c *gin.Context) {
	//我需要调用外部的py的http接口来实现图像文字识别
	// 获取上传的图片文件
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片上传失败"})
		return
	}
	defer file.Close()

	// 读取图片内容
	imgData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片读取失败"})
		return
	}

	// 调用外部 Python OCR 服务
	pyOcrUrl := "http://127.0.0.1:8080/ocr" // 假设你的py服务在这个地址
	req, err := http.NewRequest("POST", pyOcrUrl, bytes.NewReader(imgData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求OCR服务失败"})
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "调用OCR服务失败"})
		return
	}
	defer resp.Body.Close()

	ocrResult, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取OCR结果失败"})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(ocrResult, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析OCR结果失败"})
		return
	}
	fmt.Println("OCR结果:", result)
	type out struct {
		ItemName  string `json:"item_name"`
		ItemValue string `json:"item_value"`
	}
	var outResult []out
	for key, value := range result {
		// 将key转换为uint类型
		// fmt.Println("key:", key)
		// itemID, err := strconv.ParseUint(key, 10, 32)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		// 	return
		// }

		outResult = append(outResult, out{
			ItemName:  key,
			ItemValue: value.(string),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result": outResult,
	})
}
