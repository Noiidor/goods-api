package services

import (
	"goods-api/internal/broker"
	"goods-api/internal/cache"
	"goods-api/internal/dto"
	"goods-api/internal/errors"
	"goods-api/internal/models"
	"goods-api/internal/repos"
)

type GoodsService interface {
	Create(dto.CreateForm, dto.CreatePayload) (models.Good, error)
	Update(dto.UpdateForm, dto.UpdatePayload) (models.Good, error)
	Delete(dto.DeleteForm) (models.Good, error)
	GetAll(dto.GetAllForm) ([]models.Good, error)
	Reprioritize(dto.PriorityForm, dto.PriorityPayload) ([]models.Good, error)
}

type goodsService struct {
	goods    repos.GoodsRepo
	projects repos.ProjectsRepo
	cache    cache.GoodsCache
	broker   broker.MessageBroker
}

func NewGoodsService(goods repos.GoodsRepo, projects repos.ProjectsRepo, cache cache.GoodsCache, broker broker.MessageBroker) GoodsService {
	return &goodsService{
		goods:    goods,
		projects: projects,
		cache:    cache,
		broker:   broker,
	}
}

func (s goodsService) Create(form dto.CreateForm, payload dto.CreatePayload) (models.Good, error) {
	exists, err := s.projects.Exists(form.Project_id)
	if err != nil {
		return models.Good{}, err
	}

	if !exists {
		return models.Good{}, &errors.ProjectNotFoundError{Project_id: form.Project_id}
	}

	good, err := s.goods.Insert(form.Project_id, payload.Name)
	if err != nil {
		return models.Good{}, err
	}

	go s.broker.SendGood(good)

	go s.cache.AddMember(good)

	return good, nil
}

func (s goodsService) Update(form dto.UpdateForm, payload dto.UpdatePayload) (models.Good, error) {
	exist, err := s.goods.Exists(form.ID, form.Project_id)
	if err != nil {
		return models.Good{}, err
	}
	if !exist {
		return models.Good{}, &errors.GoodNotFoundError{
			ID:         form.ID,
			Project_id: form.Project_id,
		}
	}

	updateGood := models.Good{
		Name:        payload.Name,
		Description: payload.Description,
	}

	good, err := s.goods.Update(form.ID, form.Project_id, updateGood)
	if err != nil {
		return models.Good{}, err
	}

	go s.broker.SendGood(good)

	go s.cache.UpdateMember(good)

	return good, nil
}

func (s goodsService) Delete(form dto.DeleteForm) (models.Good, error) {
	exist, err := s.goods.Exists(form.ID, form.Project_id)
	if err != nil {
		return models.Good{}, err
	}
	if !exist {
		return models.Good{}, &errors.GoodNotFoundError{
			ID:         form.ID,
			Project_id: form.Project_id,
		}
	}

	good, err := s.goods.Delete(form.ID, form.Project_id)
	if err != nil {
		return models.Good{}, err
	}

	go s.broker.SendGood(good)

	go s.cache.UpdateMember(good)

	return good, nil
}

func (s goodsService) GetAll(form dto.GetAllForm) ([]models.Good, error) {
	goods := s.cache.TryGetSetMembers(form.Offset, form.Limit)
	if goods != nil {
		return goods, nil
	}

	goods, err := s.goods.GetListOffset(form.Limit, form.Offset)
	if err != nil {
		return nil, err
	}

	go func() {
		allGoods, _ := s.goods.GetAll() // Просто и эффективно для небольшого количества строк(<1кк)
		s.cache.CacheSet(allGoods)
	}()

	return goods, nil
}

func (s goodsService) Reprioritize(form dto.PriorityForm, payload dto.PriorityPayload) ([]models.Good, error) {
	exist, err := s.goods.Exists(form.ID, form.Project_id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, &errors.GoodNotFoundError{
			ID:         form.ID,
			Project_id: form.Project_id,
		}
	}

	goods, err := s.goods.Reprioritize(form.ID, form.Project_id, payload.New_priority)
	if err != nil {
		return nil, err
	}

	go s.broker.SendGoods(goods)

	go s.cache.DeleteSet()

	return goods, nil
}
