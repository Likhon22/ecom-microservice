package cartService

import (
	client "cart_service/internal/clients/product"
	"cart_service/internal/domain"
	cartRepo "cart_service/internal/repo/cart"
	cartpb "cart_service/proto/gen"
	"cart_service/utils"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type service struct {
	repo          cartRepo.Repo
	productClient client.Client
}

type Service interface {
	AddToCart(ctx context.Context, email string, req *cartpb.AddToCartRequest) (*cartpb.CartResponse, error)
}

func NewService(repo cartRepo.Repo, productClient client.Client) Service {

	return &service{
		repo:          repo,
		productClient: productClient,
	}
}

func (s *service) AddToCart(ctx context.Context, email string, req *cartpb.AddToCartRequest) (*cartpb.CartResponse, error) {

	if email == "" {
		return nil, errors.New("Unauthorized")

	}
	existingCart, err := s.repo.GetCart(ctx, email)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if existingCart == nil {
		existingCart = &domain.Cart{
			Email:     email,
			Items:     []domain.CartItem{},
			CreatedAt: time.Now(),
		}
	}

	itemIndex := utils.FindCartIndex(existingCart.Items, req.ProductId)
	if itemIndex >= 0 {
		utils.UpdateItemQuantity(&existingCart.Items[itemIndex], req.Quantity)
	} else {
		product, err := s.productClient.GetProductById(ctx, &cartpb.GetProductByIdRequest{
			Category:  req.Category,
			ProductId: req.ProductId,
		})
		if err != nil {
			return nil, err
		}

		newItem := utils.CreateCartItem(product, req.Quantity)
		existingCart.Items = append(existingCart.Items, newItem)

	}

	utils.RecalculateSubTotal(existingCart)
	existingCart.UpdatedAt = time.Now().UTC()
	savedCart, err := s.repo.AddToCart(ctx, email, existingCart)

	if err != nil {
		return nil, err
	}

	return utils.DomainCartToProto(savedCart), nil
}
