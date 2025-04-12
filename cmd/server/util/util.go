package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/DislikesSchool/EduPage2-server/config"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/redis/go-redis/v9"
)

type DataStorageConfig struct {
	Enabled     bool `json:"enabled"`
	Credentials bool `json:"credentials"`
	Messages    bool `json:"messages"`
	Timeline    bool `json:"timeline"`
}

func TTLFromType(ttlType string) time.Duration {
	switch ttlType {
	case "timeline":
		return time.Duration(config.AppConfig.Redis.TTL.Timeline) * time.Second
	case "timetable":
		return time.Duration(config.AppConfig.Redis.TTL.Timetable) * time.Second
	case "results":
		return time.Duration(config.AppConfig.Redis.TTL.Results) * time.Second
	case "dbi":
		return time.Duration(config.AppConfig.Redis.TTL.DBI) * time.Second
	default:
		return 0
	}
}

func CacheKeyFromEPClient(client *edupage.EdupageClient, key string) (string, error) {
	user, err := client.GetUser(false)
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}
	return fmt.Sprintf("%s:%s:%s", strings.Split(client.Credentials.Server, ".")[0], user.UserRow.UserID, key), nil
}

func CacheData(key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data to JSON: %w", err)
	}
	err = Rdb.Set(Ctx, key, jsonData, ttl).Err()
	if err != nil {
		return fmt.Errorf("error setting data in Redis: %w", err)
	}
	return nil
}

func IsCached(key string) (bool, error) {
	exists, err := Rdb.Exists(Ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error checking cache existence: %w", err)
	}
	return exists == 1, nil
}

func ReadCache(key string, target interface{}) (bool, error) {
	jsonData, err := Rdb.Get(Ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error getting data from Redis: %w", err)
	}

	err = json.Unmarshal([]byte(jsonData), target)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling JSON data: %w", err)
	}
	return true, nil
}
