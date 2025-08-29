package handlers

import (
	"assesment-2/database"
	"assesment-2/models"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type TradeHandler struct {
	database.ITradeDB
}

type ITradeHandler interface {
	CreateTrade(c *fiber.Ctx) error
	GetTradeBy(c *fiber.Ctx) error
	//GetPositions(c *fiber.Ctx) error
}

func NewTradeHandler(itradedb database.ITradeDB) ITradeHandler {
	return &TradeHandler{itradedb}
}

func (uh *TradeHandler) CreateTrade(c *fiber.Ctx) error {
	trade := new(models.Trade)
	err := c.BodyParser(trade)
	if err != nil {
		return err
	}

	err = trade.Validate()
	if err != nil {
		return err
	}

	trade.Status = "open"
	trade.LastTradedTime = time.Now().Unix()

	trade, err = uh.Create(trade)
	if err != nil {
		return err
	}

	return c.JSON(trade)
}

func (uh *TradeHandler) GetTradeBy(c *fiber.Ctx) error {
	orderid := c.Params("id")

	_orderid, err := strconv.Atoi(orderid)
	if err != nil {
		return errors.New("invalid Order id")
	}

	trade, err := uh.GetBy(uint(_orderid))
	if err != nil {
		log.Err(err).Msg("data might not be available or some sql issue")
		return errors.New("something went wrong or no data available with that id")
	}

	return c.JSON(trade)
}

func (uh *TradeHandler) CreatePosition(c *fiber.Ctx) error {
	trade := new(models.Trade)
	err := c.BodyParser(trade)
	if err != nil {
		return err
	}

	err = trade.Validate()
	if err != nil {
		return err
	}

	// trade.Status = "open"
	// order.LastModified = time.Now().Unix()

	trade, err = uh.ITradeDB.GetBy(trade.OrderId)
	if err != nil {
		// log here
		return fiber.NewError(fiber.StatusBadRequest, "invalid order request")
	}
	return c.JSON(trade)
}
