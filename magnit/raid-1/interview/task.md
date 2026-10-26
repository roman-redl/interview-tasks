package main

import (
	"context"
)

// Реализовать палиндром
// ---- internal/gateway/grpc/basket/add_item.go

/*
 * Реализовать метод grpc-сервера, который:
 * - добавляет товар в корзину
 * - оформляет корзину
 *
 * В оформленные корзины (basket.Status=="ordered") изменения вносить нельзя.
 * При изменении состава корзины надо пересчитывать basket.Total=sum(count*price)
 * Все элементы в корзине должны быть уникальны по ключу ProductID
 *
 * Для оформления корзины необходимо:
 * - сменить ее статус на ordered
 * - сообщить возможным потребителям с помощью сообщений в брокер.
 *
 * Дать ответ на вопрос: нужна ли здесь многопоточная защита?
 */
func (bs *BasketServer) AddItemAndOrder(ctx context.Context, req *AddItemRequest) (*EmptyResponse, error) {
	basket, err := bs.repo.Load(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

    <-- кодить тут

	return &EmptyResponse{}, nil
}

// ---- internal/service/entity/basket/basket.go
type Basket struct {
	ID     uint64
	UserID uint64
	Items  []BasketItem
	Status string
	Total  uint64
}

// ---- internal/service/entity/basket/item.go
type BasketItem struct {
	BasketID  uint64
	ProductID uint64
	Count     uint64
	Price     uint64
}

// ---- internal/gateway/grpc/basket/dependencies.go
type BasketRepository interface {
	Load(ctx context.Context, userID uint64) (*Basket, error)
	Save(ctx context.Context, b *Basket) error
}

type CheckoutProducer interface {
	SendMessage(ctx context.Context, basket Basket) error
}

// ---- internal/gateway/grpc/basket/server.go
type BasketServer struct {
	repo     BasketRepository
	producer CheckoutProducer
}

// ---- pkg/server/grpc/basket.grpc.pb.go
type BasketServiceServer interface {
	AddItemAndOrder(context.Context, *AddItemRequest) (*EmptyResponse, error)
}

type AddItemRequest struct {
	UserID    uint64
	ProductID uint64
	Price     uint64
	Count     uint64
}

type EmptyResponse struct{}
