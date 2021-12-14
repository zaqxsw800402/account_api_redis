package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"red/errs"
	"strconv"
	"time"
)

type User struct {
	Id uint `redis:"user_id"`
}

func (d Database) SaveUserID(token string, userID int) {
	d.RC.Set(ctx, token, userID, time.Hour)
}

func (d Database) GetUserID(token string) *redis.StringCmd {
	result := d.RC.Get(ctx, token)
	return result
}

// UserTimes 設定1 min內最多查詢次數
func (d Database) UserTimes(ctx context.Context, userID int, times int) *errs.AppError {
	return d.countTimes(ctx, "User", strconv.Itoa(userID), times)
}

// countTimes 設定限制次數及時間
func (d Database) countTimes(ctx context.Context, keyPrefix string, keyId string, times int) *errs.AppError {
	// 製造專屬的key
	key := fmt.Sprintf("%s:%s:times", keyPrefix, keyId)
	// 每查詢一次，值增加一
	d.RC.Incr(ctx, key)
	// 設定過期時間
	d.RC.Expire(ctx, key, time.Minute)
	// 讀取次數，並判斷是否達到次數限制
	result := d.RC.Get(ctx, key)
	value, _ := result.Int()
	if value > times {
		return errs.TooManyTimes("Log in too many times")
	}
	return nil
}
