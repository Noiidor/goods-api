package broker

import (
	"goods-api/internal/models"

	"github.com/nats-io/nats.go"
)

type MessageBroker interface {
	SendGood(good models.Good) error
	SendGoods(goods []models.Good) error
	BindGoodsReceiver() (chan *nats.Msg, error)
	SubscribeGoods(sendTo chan models.Good)
}
