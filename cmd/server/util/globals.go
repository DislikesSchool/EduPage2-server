package util

import (
	"context"

	"github.com/DislikesSchool/EduPage2-server/config"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type ClientData struct {
	CrJobId     cron.EntryID
	Client      *edupage.EdupageClient
	DataStorage *DataStorageConfig
}

var Clients = make(map[string]*ClientData)

var Ctx = context.Background()

var Cr *cron.Cron

var Rdb redis.Client
var Db *gorm.DB
var Meili meilisearch.ServiceManager
var MeiliIndex meilisearch.IndexManager

var ShouldCache = config.AppConfig.Redis.Enabled
var ShouldStore = config.AppConfig.Database.Enabled
var ShouldSearch = config.AppConfig.Meilisearch.Enabled
