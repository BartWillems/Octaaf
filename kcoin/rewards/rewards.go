package rewards

import (
	"fmt"
	"octaaf/kcoin"
	"strconv"

	goRedis "github.com/go-redis/redis"
	"github.com/gobuffalo/nulls"
	log "github.com/sirupsen/logrus"
	kalicoin "gitlab.com/bartwillems/kalicoin/pkg/models"
)

// RedisPrefix is the prefix used for storing data in redis
const RedisPrefix = "kcoin"

// RewardUsers looks at the users that deserve a reward,
// rewards them and removes them from the list
func RewardUsers(conn *goRedis.Client, event string) {
	groups := conn.SMembers(fmt.Sprintf("%v-rewards", RedisPrefix)).Val()

	for _, group := range groups {
		groupID, err := strconv.ParseInt(group, 10, 64)

		if err != nil {
			log.Errorf("Unable to fetch rewards for group: %v, error: %v", group, err)
			continue
		}

		log.Debugf("Rewarding users for group %d", groupID)

		users := getRewardUsers(conn, groupID, event)

		for _, user := range users {
			transaction, err := rewardUser(conn, groupID, user, event)

			if err != nil {
				log.Errorf("Unable to reward user with id; %v, err: %v", user, err)
			}

			log.Debugf("Rewarding user %v", user)

			if transaction.Status != kalicoin.Succeeded {
				log.Errorf("Reward transaction error for user %v, error: %v", user, transaction.FailureReason)
			}
		}

	}
}

func getRewardUsers(conn *goRedis.Client, groupID int64, event string) []int {
	redisKey := getEventKey(groupID, event)
	users := conn.SMembers(redisKey).Val()

	var userIDs []int

	for _, user := range users {
		userID, err := strconv.Atoi(user)

		if err != nil {
			log.Errorf("Unable to convert user '%v' to an integer: %v", user, err)
			continue
		}
		userIDs = append(userIDs, userID)
	}

	conn.Del(redisKey)

	return userIDs
}

// RewardUser notifies the kalicoin wallet about a user's reward and removes that user from the pending reward list
func rewardUser(conn *goRedis.Client, groupID int64, userID int, event string) (kalicoin.Transaction, error) {
	reward := kalicoin.RewardTransaction{
		Cause:    nulls.NewString(event),
		GroupID:  groupID,
		Receiver: userID,
	}

	transaction, err := kcoin.CreateTransaction(reward, "rewards", nil)

	if err != nil {
		log.Errorf("Error: %v", err)
		return transaction, err
	}

	return transaction, nil
}

// StoreUser stores the current user in redis with his reward type
// It also keeps track of the groups that participate in these kind of events
func StoreUser(conn *goRedis.Client, groupID int64, userID int, event string) {
	conn.SAdd(getEventKey(groupID, event), userID)
	storeGroup(conn, groupID)
}

func storeGroup(conn *goRedis.Client, groupID int64) {
	conn.SAdd(fmt.Sprintf("%v-rewards", RedisPrefix), groupID)
}

func getEventKey(groupID int64, event string) string {
	return fmt.Sprintf("%v-%v-%v", RedisPrefix, groupID, event)
}
