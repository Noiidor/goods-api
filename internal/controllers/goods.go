package controllers

import (
	"goods-api/internal/dto"
	"goods-api/internal/errors"
	"goods-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type GoodsController struct {
	service services.GoodsService
}

func NewGoodsController(service services.GoodsService) *GoodsController {
	return &GoodsController{service: service}
}

func (c GoodsController) Create(ctx *gin.Context) {
	var form dto.CreateForm
	if err := ctx.BindQuery(&form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}
	if err := validate.Struct(form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}

	var payload dto.CreatePayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}

	good, err := c.service.Create(form, payload)
	if err != nil {
		if _, ok := err.(*errors.ProjectNotFoundError); ok {
			errors.JsonError(ctx, 404, err)
			return
		}
		errors.JsonError(ctx, 500, err)
		return
	}

	var response dto.CreateResponse
	response.Good = good

	ctx.JSON(201, response)
}

func (c GoodsController) Delete(ctx *gin.Context) {
	var form dto.DeleteForm

	if err := ctx.BindQuery(&form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}
	if err := validate.Struct(form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}

	good, err := c.service.Delete(form)
	if err != nil {
		if _, ok := err.(*errors.GoodNotFoundError); ok {
			errors.JsonError(ctx, 404, err)
			return
		}
		errors.JsonError(ctx, 500, err)
		return
	}

	var response dto.DeleteResponse
	response.ID = good.ID
	response.Project_id = good.Project_id
	response.Removed = good.Removed

	ctx.JSON(200, response)
}
func (c GoodsController) GetAll(ctx *gin.Context) {
	var form dto.GetAllForm

	if err := ctx.BindQuery(&form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}
	if err := validate.Struct(form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}

	goods, err := c.service.GetAll(form)
	if err != nil {
		errors.JsonError(ctx, 500, err)
		return
	}

	var response dto.GetAllResponse
	response.Goods = goods
	response.Meta.Limit = form.Limit
	response.Meta.Offset = form.Offset
	response.Meta.Total = len(goods)

	removedCount := 0
	for _, g := range goods {
		if g.Removed {
			removedCount++
		}
	}

	response.Meta.Removed = removedCount

	ctx.JSON(200, response)
}
func (c GoodsController) Update(ctx *gin.Context) {
	var form dto.UpdateForm

	if err := ctx.BindQuery(&form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}
	if err := validate.Struct(form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}

	var payload dto.UpdatePayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}

	good, err := c.service.Update(form, payload)
	if err != nil {
		if _, ok := err.(*errors.GoodNotFoundError); ok {
			errors.JsonError(ctx, 404, err)
			return
		}
		errors.JsonError(ctx, 500, err)
		return
	}

	var response dto.UpdateResponse
	response.Good = good

	ctx.JSON(200, response)

}
func (c GoodsController) Reprioritize(ctx *gin.Context) {
	var form dto.PriorityForm

	if err := ctx.BindQuery(&form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}
	if err := validate.Struct(form); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}

	var payload dto.PriorityPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		errors.JsonError(ctx, 400, err)
		return
	}

	goods, err := c.service.Reprioritize(form, payload)
	if err != nil {
		if _, ok := err.(*errors.GoodNotFoundError); ok {
			errors.JsonError(ctx, 404, err)
			return
		}
		errors.JsonError(ctx, 500, err)
		return
	}

	var response dto.PriorityResponse

	priorities := make([]dto.Priority, 0, len(goods))
	for _, g := range goods {
		priority := dto.Priority{ID: g.ID, Priority: g.Priority}
		priorities = append(priorities, priority)
	}

	response.Priorities = priorities

	ctx.JSON(200, response)
}
