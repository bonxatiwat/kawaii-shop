package ordersHandlers

import (
	"strings"
	"time"

	"github.com/bonxatiwat/kawaii-shop-tutortial/config"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/entities"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/orders"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ordersHandlersErrCode string

const (
	findOneOrderErr ordersHandlersErrCode = "orders-001"
	findOrderErr    ordersHandlersErrCode = "orders-002"
	insertOrderErr  ordersHandlersErrCode = "orders-003"
	updateOrderErr  ordersHandlersErrCode = "orders-004"
)

type IOrdersHandler interface {
	FindOneOrder(c *fiber.Ctx) error
	FindOrder(c *fiber.Ctx) error
	InsertOrder(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
}

type ordersHandler struct {
	cfg           config.IConfig
	ordersUsecase ordersUsecases.IOrdersUsecase
}

func OrdersHandler(cfg config.IConfig, ordersUsecase ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:           cfg,
		ordersUsecase: ordersUsecase,
	}
}

func (h *ordersHandler) FindOneOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")
	order, err := h.ordersUsecase.FindOneOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneOrderErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}

func (h *ordersHandler) FindOrder(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(findOrderErr),
			err.Error(),
		).Res()
	}

	// pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	// Sort
	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}
	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}

	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}

	// Date YYYY-MM-DD
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadGateway.Code,
				string(findOrderErr),
				"start date in invalid",
			).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}

	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadGateway.Code,
				string(findOrderErr),
				"end date in invalid",
			).Res()
		}

		req.EndDate = end.Format("2006-01-02")
	}

	// Usecase
	orders := h.ordersUsecase.FindOrder(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, orders).Res()
}

func (h *ordersHandler) InsertOrder(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	req := &orders.Order{
		Products: make([]*orders.ProductsOrder, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}
	if len(req.Products) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertOrderErr),
			"products are empty",
		).Res()
	}
	if c.Locals("userRoleId").(int) != 2 {
		req.UserId = userId
	}

	req.Status = "waiting"
	req.TotalPaid = 0

	order, err := h.ordersUsecase.InsertOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, order).Res()
}

func (h *ordersHandler) UpdateOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")
	req := new(orders.Order)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}
	req.Id = orderId
	statusMap := map[string]string{
		"waiting":   "waiting",
		"shiping":   "shiping",
		"completed": "completed",
		"canceled":  "canceled",
	}

	if c.Locals("userRoleId").(int) == 2 {
		req.Status = statusMap[strings.ToLower(req.Status)]
	} else if strings.ToLower(req.Status) == statusMap["canceled"] {
		req.Status = statusMap["canceled"]
	}

	if req.TransferSlip != nil {
		if req.TransferSlip.Id == "" {
			req.TransferSlip.Id = uuid.NewString()
		}
		if req.TransferSlip.CreatedAt == "" {
			loc, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				return entities.NewResponse(c).Error(
					fiber.ErrInternalServerError.Code,
					string(updateOrderErr),
					err.Error(),
				).Res()
			}
			now := time.Now().In(loc)

			// YYYY-MM-DD HH:MM:SS
			// 2006-01-02 15:04:05
			req.TransferSlip.CreatedAt = now.Format("2006-01-02 15:04:05")
		}
	}

	order, err := h.ordersUsecase.UpdateOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
