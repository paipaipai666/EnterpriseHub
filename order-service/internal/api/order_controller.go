package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/domain/dto"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/service"
)

type OrderController interface {
	CreateOrder(ctx *gin.Context)
	PayForOrder(ctx *gin.Context)
	GetOrderById(ctx *gin.Context)
	GetOrderList(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
}

type orderControllerImpl struct {
	service service.OrderService
}

func NewOrderController(service service.OrderService) OrderController {
	return &orderControllerImpl{
		service: service,
	}
}

func (oci *orderControllerImpl) CreateOrder(ctx *gin.Context) {
	orderParam := &dto.OrderDTO{}

	err := ctx.ShouldBindJSON(&orderParam)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed",
			"data":    err,
		})
		return
	}

	orderId, err := oci.service.CreateOrder(orderParam.UserId, orderParam.Amount)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed",
			"data":    err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    orderId,
	})
}

func (oci *orderControllerImpl) PayForOrder(ctx *gin.Context) {
	method := &dto.PaymentDTO{}
	orderId := ctx.Param("id")

	err := ctx.ShouldBindJSON(&method)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed1",
			"data":    err.Error(),
		})
		return
	}

	fmt.Println(method.Method.ToMethod())
	id, err := oci.service.Pay(ctx, method.Method.ToMethod(), orderId)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    id,
	})
}

func (oci *orderControllerImpl) GetOrderById(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := oci.service.GetOrder(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed",
			"data":    err,
		})
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    order,
	})
}

func (oci *orderControllerImpl) GetOrderList(ctx *gin.Context) {
	userId := ctx.Param("user_id")

	orders, err := oci.service.GetOrderList(userId)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed",
			"data":    err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    orders,
	})
}

func (oci *orderControllerImpl) CancelOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	err := oci.service.Cancel(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed",
			"data":    err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    "订单 " + id + " 已取消！",
	})
}
