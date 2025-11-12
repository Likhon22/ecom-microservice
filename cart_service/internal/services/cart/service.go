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
	GetCart(ctx context.Context, email string) (*cartpb.CartResponse, error)
	UpdateCart(ctx context.Context, email string, req *cartpb.UpdateCartItemRequest) (*cartpb.CartResponse, error)
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

func (s *service) GetCart(ctx context.Context, email string) (*cartpb.CartResponse, error) {

	if email == "" {
		return nil, errors.New("Unauthorized")

	}
	resp, err := s.repo.GetCart(ctx, email)
	if err != nil {
		return nil, err

	}
	return utils.DomainCartToProto(resp), nil

}

func (s *service) UpdateCart(ctx context.Context, email string, req *cartpb.UpdateCartItemRequest) (*cartpb.CartResponse, error) {

	if email == "" {
		return nil, errors.New("Unauthorized")
	}
	cart, err := s.repo.GetCart(ctx, email)
	if err != nil {
		return nil, err
	}
	itemIndex := utils.FindCartIndex(cart.Items, req.ProductId)
	if itemIndex < 0 {
		return nil, errors.New("product not found in cart")
	}
	if req.Quantity == 0 {
		cart.Items = append(cart.Items[:itemIndex], cart.Items[itemIndex+1:]...)

	} else {
		cart.Items[itemIndex].Quantity = req.Quantity
		cart.Items[itemIndex].Subtotal = cart.Items[itemIndex].Price * float64(req.Quantity)
	}
	utils.RecalculateSubTotal(cart)
	cart.UpdatedAt = time.Now().UTC()
	savedCart, err := s.repo.AddToCart(ctx, email, cart)
	if err != nil {
		return nil, err
	}
	return utils.DomainCartToProto(savedCart), nil
}
