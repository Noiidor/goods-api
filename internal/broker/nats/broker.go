package nats

import (
	"bytes"
	"encoding/gob"
	"goods-api/internal/broker"
	"goods-api/internal/errors"
	"goods-api/internal/models"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	GOODS_CHANNEL = "goods"
)

type messageBroker struct {
	nats *nats.Conn
}

func New(conn *nats.Conn) broker.MessageBroker {
	return &messageBroker{nats: conn}
}

func (b messageBroker) SendGood(good models.Good) error {
	good.Created_at = time.Now() // EventTime в CH, костыль что бы не создавать отдельную модель с таким же полем, но другим названием
	var encodedGood bytes.Buffer
	enc := gob.NewEncoder(&encodedGood)
	err := enc.Encode(good)
	if err != nil {
		errors.AsyncErrors <- err
	}

	err = b.nats.Publish(GOODS_CHANNEL, encodedGood.Bytes())
	if err != nil {
		errors.AsyncErrors <- err
	}
	return nil
}

func (b messageBroker) SendGoods(goods []models.Good) error {
	now := time.Now()
	for _, v := range goods {
		v.Created_at = now
		go b.SendGood(v)
	}

	return nil
}

func (b messageBroker) SubscribeGoods(sendTo chan models.Good) {
	b.nats.Subscribe(GOODS_CHANNEL, func(msg *nats.Msg) {
		var buffer bytes.Buffer
		buffer.Write(msg.Data)
		var decodedGood models.Good

		decoder := gob.NewDecoder(&buffer)
		err := decoder.Decode(&decodedGood)
		if err != nil {
			errors.AsyncErrors <- err
		}
		sendTo <- decodedGood
	})
}

func (b messageBroker) BindGoodsReceiver() (chan *nats.Msg, error) {
	ch := make(chan *nats.Msg)
	_, err := b.nats.ChanSubscribe(GOODS_CHANNEL, ch)
	if err != nil {
		return nil, err
	}

	return ch, nil

}
