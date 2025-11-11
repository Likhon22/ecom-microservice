package cartService

import (
	client "cart_service/internal/clients/product"
	"cart_service/internal/domain"
	cartRepo "cart_service/internal/repo/cart"
	cartpb "cart_service/proto/gen"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	itemsExists := false
	for i, item := range existingCart.Items {
		if item.ProductID == req.ProductId {
			existingCart.Items[i].Quantity += req.Quantity
			existingCart.Items[i].Subtotal = existingCart.Items[i].Price * float64(existingCart.Items[i].Quantity)
			itemsExists = true
			break

		}

	}
	var product *cartpb.Product
	if !itemsExists {
		product, err = s.productClient.GetProductById(ctx, &cartpb.GetProductByIdRequest{
			Category:  req.Category,
			ProductId: req.ProductId,
		})
		if err != nil {
			return nil, err
		}

		newItem := domain.CartItem{
			ProductID:   req.ProductId,
			Category:    req.Category,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    req.Quantity,
			ImageURL:    product.ImageUrls[0],
			Subtotal:    product.Price * float64(req.Quantity),
		}
		existingCart.Items = append(existingCart.Items, newItem)
	}
	existingCart.TotalItems = 0
	existingCart.Subtotal = 0
	for _, item := range existingCart.Items {
		existingCart.TotalItems += item.Quantity
		existingCart.Subtotal += item.Subtotal

	}
	existingCart.UpdatedAt = time.Now().UTC()
	cart, err := s.repo.AddToCart(ctx, email, existingCart)

	if err != nil {
		return nil, err
	}
	pbItems := make([]*cartpb.CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		pbItems = append(pbItems, &cartpb.CartItem{
			ProductId:   item.ProductID,
			Category:    item.Category,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
			ImageUrl:    item.ImageURL,
			Subtotal:    item.Subtotal,
		})
	}
	return &cartpb.CartResponse{
		Email:      email,
		Items:      pbItems,
		TotalItems: cart.TotalItems,
		Subtotal:   cart.Subtotal,
		CreatedAt:  timestamppb.New(cart.CreatedAt),
		UpdatedAt:  timestamppb.New(cart.UpdatedAt),
	}, nil
}
