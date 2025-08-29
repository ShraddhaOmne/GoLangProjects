package handlers

import (
	"errors"
	"food-delivery/database"
	"food-delivery/messaging"
	"food-delivery/models"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderHandler struct {
	database.IOrderDB // prmoted field
	database.IOrderEventDB
}

// GetOrderBy implements IOrderHandler.
func (o *OrderHandler) GetOrderBy(c *fiber.Ctx) error {
	panic("unimplemented")
}

type IOrderHandler interface {
	CreateOrder(msg *messaging.Messaging) func(c *fiber.Ctx) error
	GetOrder() func(c *fiber.Ctx) error
}

func NewOrderHandler(iorderdb database.IOrderDB, iordereventdb database.IOrderEventDB) IOrderHandler {
	return &OrderHandler{
		IOrderDB:      iorderdb,
		IOrderEventDB: iordereventdb,
	}
}

func (o *OrderHandler) CreateOrder(msg *messaging.Messaging) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		order := new(models.Orders)
		err := c.BodyParser(order)
		if err != nil {
			return err
		}
		err = order.Validate()
		if err != nil {
			return err
		}

		order.Status = "PLACED"
		order.CreatedAt = time.Now().Unix()
		order.UpdatedAt = time.Now().Unix()

		order, err = o.Create(order)
		if err != nil {
			return err
		}
		if err := o.IOrderEventDB.CreateEvent(order.OrderId, order.Status); err != nil {
			log.Printf("failed to insert initial event for order %d: %v", order.OrderId, err)
		}
		msg.ChMessaging <- order.ToBytes()
		go o.simulateOrderProgression(order.OrderId, order.Status)

		return c.JSON(order)
	}
}
func (o *OrderHandler) GetOrder() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// parse id
		idParam := c.Params("order_id")
		orderID, err := strconv.Atoi(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid order_id",
			})
		}

		// fetch from DB
		order, err := o.IOrderDB.GetWithEvents(uint(orderID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "order not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// return combined response
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"order":  order,
			"events": order.Events,
		})
	}
}

func (o *OrderHandler) simulateOrderProgression(orderID uint, status string) {
	statuses := []string{"PREPARING", "COOKING", "OUT_FOR_DELIVERY", "DELIVERED"}

	for _, status := range statuses {
		time.Sleep(time.Second * 3)

		// update order
		_ = o.IOrderDB.UpdateStatus(orderID, status)
		if err := o.IOrderEventDB.CreateEvent(orderID, status); err != nil {
			log.Printf("failed to insert order_event for order %d: %v", orderID, err)
		}

	}

}
