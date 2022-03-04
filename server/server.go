/*
Copyright © 2022 shfz

*/
package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shfz/shfz/fuzz"
	"github.com/shfz/shfz/model"
)

func Server() {
	r := gin.Default()

	// ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// filibからFuzzを取得する
	r.POST("/fuzz", func(c *gin.Context) {
		var param model.FuzzInfo
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		text, err := fuzz.GetFuzz(param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"fuzz": text,
		})
	})

	// APIとFuzzの対応関係を保存
	r.POST("/api", func(c *gin.Context) {
		var param model.ApiParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		if err := fuzz.SetFuzzTexts(param); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": "ok",
		})
	})

	// filibからclientエラー情報を送る
	r.POST("/client/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no id",
			})
			return
		}
		var param model.ClientFeedback
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		if err := fuzz.SetClientFeedback(id, param); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": "ok",
		})
	})

	// webアプリからframeとserverエラー情報を送る
	r.POST("/server/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no id",
			})
			return
		}
		var param model.ServerFeedback
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		if err := fuzz.SetServerFeedback(id, param); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": "ok",
		})
	})

	r.GET("/data", func(c *gin.Context) {
		stat := fuzz.GetApiData()
		c.JSON(http.StatusOK, gin.H{
			"status": stat,
		})
	})

	r.POST("/report", func(c *gin.Context) {
		var param model.ReportReq
		if err := c.Bind(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		str, err := fuzz.GenReport(param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		c.String(200, str)
	})

	if err := r.Run(":53653"); err != nil {
		panic(err)
	}
}
