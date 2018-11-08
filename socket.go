package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type OctaafSocket struct {
	*gin.Engine
}

func NewOctaafSocket() *OctaafSocket {
	router := gin.Default()
	o := &OctaafSocket{router}
	api := o.Group("/api/v1")
	{
		api.POST("/reload", o.Reload)
	}

	return o
}

func (o *OctaafSocket) Listen() error {
	return o.Run("127.0.0.1:8127")
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
	}

	c.JSON(status, result)
}

// SocketWriter writes messages to the local octaaf socket
func SocketWriter() {
	log.Fatal("To be implemented")
}
