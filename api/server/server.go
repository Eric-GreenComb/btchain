package server

import (
	"github.com/axengine/btchain/api/config"
	"github.com/axengine/btchain/api/handlers"
	"go.uber.org/zap"
	"os"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg     *config.Config
	handler *handlers.Handler
}

func NewServer(logger *zap.Logger, cfg *config.Config) *Server {
	handler := handlers.NewHandler(logger, cfg)

	return &Server{
		cfg:     cfg,
		handler: handler,
	}
}

func (s *Server) Start() {
	router := gin.Default()

	router.Handle("HEAD", "/", func(context *gin.Context) {
		context.String(200, "%s", "success")
	})

	v1 := router.Group("/v1")
	{
		v1.POST("/transactions", s.handler.SendTransactions)
		//v1.GET("/nonce/:address", s.handler.QueryNonce)
		v1.GET("/genkey", s.handler.GenKey)
		v1.GET("/accounts/:address", s.handler.QueryAccount)
		v1.GET("/transactions/:txhash", s.handler.QuerySingleTx)
		v1.GET("/transactions", s.handler.QueryTxs)
		v1.GET("/accounts/:address/transactions", s.handler.QueryAccTxs)
	}

	if len(os.Args) > 1 && os.Args[1] == "version" {
		return
	}
	// s.handler.ReqServerInfo()
	router.Run(s.cfg.Bind)
}
