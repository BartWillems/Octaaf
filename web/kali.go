package web

import (
	"octaaf/models"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

func quote(c *gin.Context) {
	quote := models.Quote{}
	err := conn.Postgres.Where("chat_id = ?", conn.KaliID).Order("random()").Limit(1).First(&quote)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	config := tgbotapi.ChatConfigWithUser{
		ChatID:             conn.KaliID,
		SuperGroupUsername: "",
		UserID:             quote.UserID}

	user, err := conn.Octaaf.GetChatMember(config)

	quoteMap := structs.Map(quote)

	if err != nil {
		quoteMap["from"] = err
	} else {
		quoteMap["from"] = user
	}

	c.JSON(200, quoteMap)
}
