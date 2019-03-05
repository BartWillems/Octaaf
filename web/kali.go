package web

import (
	"fmt"
	"octaaf/models"
	"octaaf/trump"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func getQuote(filter string) (models.Quote, string, error) {
	quote := models.Quote{}
	err := quote.Search(conn.Postgres, conn.KaliID, filter)

	if err != nil {
		return models.Quote{}, "", err
	}

	config := tgbotapi.ChatConfigWithUser{
		ChatID:             conn.KaliID,
		SuperGroupUsername: "",
		UserID:             quote.UserID}

	user, err := conn.Octaaf.GetChatMember(config)

	if err != nil {
		return models.Quote{}, "", err
	}

	return quote, user.User.UserName, nil

}

func quote(c *gin.Context) {
	quote, username, err := getQuote(c.Query("filter"))

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	quoteMap := structs.Map(quote)
	quoteMap["from"] = username

	c.JSON(200, quoteMap)
}

func presidentialQuote(c *gin.Context) {
	quote, username, err := getQuote(c.Query("filter"))

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	msg := fmt.Sprintf(`"%v"`, quote.Quote)
	msg += fmt.Sprintf("\n    ~@%v", username)
	img, err := trump.Order(*conn.Trump, conn.TrumpCfg, msg)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Data(200, "image/png", img)
}
