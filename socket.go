package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const SOCKET_URI = "127.0.0.1:8127"
const SOCKET_PATH = "api/v1"

type OctaafSocket struct {
	*gin.Engine
}

func NewOctaafSocket() *OctaafSocket {
	router := gin.Default()

	if settings.Environment == "production" {
		gin.SetMode("release")
	}

	o := &OctaafSocket{router}
	api := o.Group(fmt.Sprintf("/%s", SOCKET_PATH))
	{
		api.POST("/reload", o.Reload)
	}

	return o
}

func (o *OctaafSocket) Listen() error {
	return o.Run(SOCKET_URI)
}

func (o *OctaafSocket) Reload(c *gin.Context) {
	_, err := settings.Load()

	status := 200
	result := gin.H{
		"message": "Settings reloaded.",
	}

	if err != nil {
		status = 500
		result["message"] = err
		log.Errorf("There was an error while reloading the settings: %v", err)
	} else {
		log.Info("Settings sucessfully reloaded.")
	}

	c.JSON(status, result)
}
