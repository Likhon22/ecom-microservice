package orderService

import (
	"order_service/internal/kafka"
	orderRepo "order_service/internal/repo/order"
)

type service struct {
	repo     orderRepo.Repo
	producer kafka.Producer
	consumer kafka.Consumer
}

type Service interface {
}

func NewService(producer kafka.Producer, consumer kafka.Consumer, repo orderRepo.Repo) Service {

	return &service{
		producer: producer,
		consumer: consumer,
		repo:     repo,
	}
}
