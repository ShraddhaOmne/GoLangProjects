package handlers

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"miniMutualFund/database"
	"miniMutualFund/messaging"
	"miniMutualFund/models"
	"miniMutualFund/redis"
	"time"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	database.IOrderDB // prmoted field
}

type IOrderHandler interface {
	CreateOrder(msg *messaging.Messaging, rdb *redis.RedisClient) func(c *fiber.Ctx) error
	GetAllOrders() func(c *fiber.Ctx) error
}

func NewOrderHandler(iorderdb database.IOrderDB) IOrderHandler {
	return &OrderHandler{
		IOrderDB: iorderdb,
	}
}

var ctx = context.Background()

func (o *OrderHandler) CreateOrder(msg *messaging.Messaging, rdb *redis.RedisClient) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		order := new(models.PlaceOrder)
		err := c.BodyParser(order)
		if err != nil {
			return err
		}

		err = order.Validate()
		if err != nil {
			return err
		}
		// ✅ Get NAV from Redis
		navKey := "nav:latest:" + order.SchemeCode
		navVal, err := rdb.Client.Get(ctx, navKey).Float64()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to get NAV from redis"})
		}

		order.NAVUsed = navVal
		order.Units = order.Amount / navVal
		order.Units = math.Round(order.Units*10000) / 10000

		order.Status = "PLACED"
		order.CreatedAt = time.Now().Unix()
		msg.ChMessaging <- order.ToBytes()
		order, err = o.Create(order)
		if err != nil {
			return err
		}
		go func(orderID uint) {
			time.Sleep(5 * time.Second)

			ord, err := o.GetByID(orderID)
			if err != nil {
				return
			}
			ord.Status = "CONFIRMED"
			//ord.ConfirmedAt = time.Now().Unix()

			_, err = o.Update(ord)
			if err != nil {
				return
			}
		}(order.OrderId)

		return c.JSON(order)
	}
}

func (o *OrderHandler) GetAllOrders() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		orders, err := o.GetAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "failed to fetch orders: " + err.Error(),
			})
		}
		return c.JSON(orders)
	}
}
func StartNAVSimulator(ctx context.Context, rdb *redis.RedisClient, schemeCodes []string) {
	go func() {
		rand.Seed(time.Now().UnixNano())
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("🛑 Stopping NAV simulator...")
				return
			default:
				for _, scheme := range schemeCodes {
					navVal := 10 + rand.Float64()*(200-10)
					navVal = math.Round(navVal*10000) / 10000
					navKey := "nav:latest:" + scheme
					err := rdb.Client.Set(ctx, navKey, navVal, 0).Err()
					if err != nil {
						fmt.Println("❌ Failed to update NAV for", scheme, ":", err)
					} else {
						fmt.Println("✅ Updated NAV for", scheme, ":", navVal)
					}
				}
				<-ticker.C
			}
		}
	}()
}
