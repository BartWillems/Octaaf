package web

import "github.com/gin-gonic/gin"

func health(c *gin.Context) {
	// HTTP status code, either 200 or 500
	status := 200
	// User facing status message, contains all services
	var statusMessage = gin.H{
		"redis":    "Ok",
		"postgres": "Ok",
		"telegram": "Ok",
	}

	// Redis
	_, redisErr := conn.Redis.Ping().Result()
	if redisErr != nil {
		statusMessage["redis"] = redisErr
		status = 500
	}

	// Postgres
	postgresErr := conn.Postgres.RawQuery("SELECT COUNT(pid) FROM pg_stat_activity;").Exec()
	if postgresErr != nil {
		statusMessage["postgres"] = postgresErr
		status = 500
	}

	// Telegram
	_, telegramErr := conn.Octaaf.GetMe()
	if telegramErr != nil {
		statusMessage["telegram"] = telegramErr
		status = 500
	}

	c.JSON(status, statusMessage)
}
