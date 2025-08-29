package handlers

import (
	"errors"
	"strconv"
	"time"
	"users-service/database"
	"users-service/messaging"
	"users-service/models"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	database.IUserDB // prmoted field
}

type IUserHandler interface {
	CreateUser(msg *messaging.Messaging) func(c *fiber.Ctx) error
	GetUserBy(c *fiber.Ctx) error
	GetUsersByLimit(c *fiber.Ctx) error
	CreateOrder(c *fiber.Ctx) error
}

func NewUserHandler(iuserdb database.IUserDB) IUserHandler {
	return &UserHandler{iuserdb}
}

func (uh *UserHandler) CreateUser(msg *messaging.Messaging) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		err := c.BodyParser(user)
		if err != nil {
			return err
		}

		err = user.Validate()
		if err != nil {
			return err
		}

		user.Status = "active"
		user.LastModified = time.Now().Unix()

		user, err = uh.Create(user)
		if err != nil {
			return err
		}
		msg.ChMessaging <- user.ToBytes()
		return c.JSON(user)
	}
}

func (uh *UserHandler) GetUserBy(c *fiber.Ctx) error {
	id := c.Params("id") // Retrieves the value of ":id"

	_id, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid id")
	}

	user, err := uh.GetBy(uint(_id))
	if err != nil {
		log.Err(err).Msg("data might not be available or some sql issue")
		return errors.New("something went wrong or no data available with that id")
	}

	return c.JSON(user)
}

func (uh *UserHandler) GetUsersByLimit(c *fiber.Ctx) error {
	limit := c.Params("limit") // Retrieves the value of ":id"

	l, err := strconv.Atoi(limit)
	if err != nil {
		return errors.New("invalid limit")
	}

	offset := c.Params("offset") // Retrieves the value of ":id"

	of, err := strconv.Atoi(offset)
	if err != nil {
		return errors.New("invalid offset")
	}

	users, err := uh.GetByLimit(l, of)
	if err != nil {
		log.Err(err).Msg("data might not be available or some sql issue")
		return errors.New("something went wrong or no data available with that id")
	}

	return c.JSON(users)
}

func (uh *UserHandler) CreateOrder(c *fiber.Ctx) error {
	order := new(models.Order)
	err := c.BodyParser(order)
	if err != nil {
		return err
	}

	err = order.Validate()
	if err != nil {
		return err
	}

	order.Status = "active"
	order.LastModified = time.Now().Unix()

	order, err = uh.IUserDB.CreateOrder(order)
	if err != nil {
		// log here
		return fiber.NewError(fiber.StatusBadRequest, "invalid order request")
	}
	return c.JSON(order)
}
