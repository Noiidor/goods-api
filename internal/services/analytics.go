package services

import (
	"context"
	"goods-api/internal/broker"
	"goods-api/internal/errors"
	"goods-api/internal/models"
	"goods-api/internal/repos"
)

const (
	GOODS_BATCH_SIZE = 5 // Размер пачки для отправки в ClickHouse
)

type AnalyticsService interface {
	HandleQueue(ctx context.Context)
}

type analyticsService struct {
	repo   repos.GoodsAnalytics
	broker broker.MessageBroker
}

func NewAnalyticsService(repo repos.GoodsAnalytics, broker broker.MessageBroker) AnalyticsService {
	return &analyticsService{repo: repo, broker: broker}
}

// Объекты на логгирование в сервис приходят через NATS от самого же себя
// Фактически конечно лог через очередь NATS это должен быть отдельный микросервис
func (s analyticsService) HandleQueue(ctx context.Context) {
	goodsCh := make(chan models.Good)
	s.broker.SubscribeGoods(goodsCh)

	go func() { // Потенциальный bottleneck на приём из канала, но можно отскейлить просто запустив несколько таких горутин
		buffer := make([]models.Good, 0)

		for {
			select {
			case <-ctx.Done():
				if err := s.repo.InsertBatch(buffer); err != nil {
					errors.AsyncErrors <- err
				}

			case good := <-goodsCh:
				buffer = append(buffer, good)

				if len(buffer) < GOODS_BATCH_SIZE {
					continue
				}

				if err := s.repo.InsertBatch(buffer); err != nil {
					errors.AsyncErrors <- err
				}
				buffer = nil
			}
		}
	}()
}
