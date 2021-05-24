package http

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strconv"
	"sync-ethereum/internal/config"
	pkgErrors "sync-ethereum/internal/errors"
	"sync-ethereum/internal/model"
	"sync-ethereum/internal/service"
	"sync-ethereum/pkg/mq"

	ginLogger "github.com/gin-contrib/logger"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	_RequestIDHeaderName = "X-Request-Id"
)

type HttpServer struct {
	config     config.Config
	httpServer *http.Server
	engine     *gin.Engine
	logger     zerolog.Logger
	mq         mq.MQ
	storageSvc service.StorageService
}

func (server *HttpServer) setRouter() {
	server.engine.Use(gin.Recovery())
	server.engine.Use(gzip.Gzip(gzip.DefaultCompression))
	server.engine.Use(ginLogger.SetLogger(ginLogger.Config{
		Logger: &server.logger,
		UTC:    true,
	}))
	server.engine.Use(server.RequestIDMiddleware)
	server.engine.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Page Not Found")
	})
	server.engine.NoMethod(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
	})

	{
		apiV1 := server.engine.Group("/api/v1")
		apiV1.GET("/blocks", server.GetBlocks)
		apiV1.GET("/blocks/:id", server.GetBlock)
		apiV1.GET("/transaction/:txhash", server.GetTransation)
	}
}

func NewHttpServer(config config.Config, logger zerolog.Logger, mq mq.MQ, storageSvc service.StorageService) *HttpServer {
	httpServer := &HttpServer{
		config:     config,
		engine:     gin.Default(),
		logger:     logger,
		mq:         mq,
		storageSvc: storageSvc,
	}
	httpServer.setRouter()

	return httpServer
}

func (server *HttpServer) Run(addr string) error {
	server.httpServer = &http.Server{
		Addr:    addr,
		Handler: server.engine,
	}
	return server.httpServer.ListenAndServe()
}

func (server *HttpServer) Shutdown() error {
	if err := server.httpServer.Close(); err != nil {
		return err
	}
	if err := server.mq.Close(); err != nil {
		return err
	}
	if err := server.storageSvc.Close(); err != nil {
		return err
	}
	return nil
}

func (server *HttpServer) RequestIDMiddleware(ctx *gin.Context) {
	uuid := uuid.New().String()
	if requestID := ctx.GetHeader(_RequestIDHeaderName); len(requestID) > 0 {
		uuid = requestID
	}
	ctx.Set(_RequestIDHeaderName, uuid)
	ctx.Header(_RequestIDHeaderName, uuid)

	ctx.Next()
}

func (server *HttpServer) GetBlocks(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	limit := 10
	convLimit, err := strconv.Atoi(limitStr)
	if err != nil {
		server.logger.Warn().Err(err).Msg("input param limit is invalid")
	} else {
		limit = convLimit
	}
	blocks, err := server.storageSvc.ListBlock(ctx, model.Block{}, model.Pagination{
		PerPage: int64(limit),
	}, model.Sorting([]model.SortField{{Field: "block_num", Order: model.SortDESC}}))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		server.logger.Error().Err(err).Msg("list block error")
		return
	}

	respBlock := make([]Block, len(blocks))
	for i, block := range blocks {
		respBlock[i] = Block{
			BlockNumber: block.BlockNumber,
			BlockHash:   block.BlockHash,
			BlockTime:   block.BlockTime,
			ParentHash:  block.ParentHash,
			IsStable:    block.IsStable,
		}
	}

	ctx.JSON(http.StatusOK, GetBlocksResponse{respBlock})
}

func (server *HttpServer) GetBlock(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		server.logger.Warn().Err(err).Msg("input param id is invalid")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	bigI := big.NewInt(id)
	block, err := server.storageSvc.GetBlock(ctx, model.Block{
		BlockNumber: model.GormBigInt(*bigI),
	})
	if err != nil {
		if errors.Is(err, pkgErrors.ErrResourceNotFound) {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		server.logger.Error().Err(err).Msg("get block error")
		return
	}

	if !block.IsStable {
		server._Compensate(ctx, block)
	}

	transactions := make([]string, len(block.Transaction))
	for i, tx := range block.Transaction {
		transactions[i] = tx.TXHash
	}

	ctx.JSON(http.StatusOK, GetBlockResponse{
		BlockNumber:  block.BlockNumber,
		BlockHash:    block.BlockHash,
		BlockTime:    block.BlockTime,
		ParentHash:   block.ParentHash,
		IsStable:     block.IsStable,
		Transactions: transactions,
	})
}

func (server *HttpServer) _Compensate(ctx context.Context, block model.Block) {
	currentBlock, err := server.storageSvc.GetCurrentBlockNumber(ctx)
	if err != nil {
		server.logger.Error().Err(err).Msg("get current block number error")
		return
	}
	if block.BlockNumber.Int64() < (currentBlock.Int64() - int64(server.config.Scheduler.UnstableNumber)) {
		message := model.CrawlerMessage{
			IsStable:    true,
			BlockNumber: block.BlockNumber,
		}
		messageBytes, err := json.Marshal(message)
		if err != nil {
			server.logger.Error().Int64("block_number", block.BlockNumber.Int64()).Err(err).Msg("marshal crawler message error")
			return
		}
		if err := server.mq.Publish(server.config.Crawler.Topic, uuid.New().String(), messageBytes); err != nil {
			server.logger.Error().Int64("block_number", block.BlockNumber.Int64()).Err(err).Msg("push crawler id error")
			return
		}
		server.logger.Info().Int64("block_number", block.BlockNumber.Int64()).Msg("compensate block")
	}
}

func (server *HttpServer) GetTransation(ctx *gin.Context) {
	txhash := ctx.Param("txhash")

	transaction, err := server.storageSvc.GetTransaction(ctx, model.Transaction{
		TXHash: txhash,
	})
	if err != nil {
		if errors.Is(err, pkgErrors.ErrResourceNotFound) {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		server.logger.Error().Err(err).Msg("get block error")
		return
	}

	logs := make([]TransactionLog, len(transaction.Logs))
	for i, log := range transaction.Logs {
		logs[i] = TransactionLog{
			Index: log.Index,
			Data:  log.Data,
		}
	}

	ctx.JSON(http.StatusOK, GetTransactionResponse{
		TXHash: transaction.TXHash,
		From:   transaction.From,
		To:     transaction.To,
		Nonce:  transaction.Nonce,
		Data:   transaction.Data,
		Value:  transaction.Value,
		Logs:   logs,
	})
}
