package handlers

import (
	"fmt"
	"github.com/axengine/btchain/api/cache"
	"github.com/axengine/btchain/api/config"
	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/rpc/client"
	"go.uber.org/zap"
	"sync"
)

type Handler struct {
	client *client.HTTP
	logger *zap.Logger
	mu     sync.Mutex
	cache  *cache.Cache
}

func NewHandler(logger *zap.Logger, cfg *config.Config) *Handler {
	var h Handler
	h.client = client.NewHTTP(cfg.RPC, "/websocket")
	h.logger = logger

	h.cache = cache.New(cfg.Cache.RedisStore, 300)
	if cfg.Cache.RedisStore {
		h.cache.InitRedisCache(cfg.Cache.Dial, cfg.Cache.Password, cfg.Cache.Db, cfg.Cache.MaxConnections)
	}
	return &h
}

func (hd *Handler) responseWrite(ctx *gin.Context, isSuccess bool, result interface{}) {
	ret := gin.H{
		"isSuccess": isSuccess,
	}

	if isSuccess {
		ret["result"] = result
	} else {
		ret["message"] = result
	}

	ctx.JSON(200, ret)

	fmt.Printf("===========raw request url: %s\n", ctx.Request.URL.String())
	fmt.Printf("===========raw response result: %v\n", result)
}
