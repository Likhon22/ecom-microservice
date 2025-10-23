package types

import userpb "auth_service/proto/gen"

type CreateCustomerInput struct {
	Name      string
	Email     string
	Password  string
	Phone     string
	Address   string
	AvatarURL string
}

type CreateCustomerResult struct {
	Name      string
	Email     string
	Role      string
	Status    string
	Phone     string
	Address   string
	AvatarURL string
	IsDeleted bool
}

type DeleteCustomerResult struct {
	Message string
}

func FromProtoCustomer(res *userpb.CreateCustomerResponse) *CreateCustomerResult {
	return &CreateCustomerResult{
		Name:      res.GetName(),
		Email:     res.GetEmail(),
		Role:      res.GetRole(),
		Status:    res.GetStatus(),
		Phone:     res.GetPhone(),
		Address:   res.GetAddress(),
		AvatarURL: res.GetAvatarUrl(),
		IsDeleted: res.GetIsDeleted(),
	}
}

func FromProtoCustomers(list []*userpb.CreateCustomerResponse) []*CreateCustomerResult {
	results := make([]*CreateCustomerResult, 0, len(list))
	for _, item := range list {
		results = append(results, FromProtoCustomer(item))
	}
	return results
}

func (r *CreateCustomerResult) ToProto() *userpb.CreateCustomerResponse {
	return &userpb.CreateCustomerResponse{
		Name:      r.Name,
		Email:     r.Email,
		Role:      r.Role,
		Status:    r.Status,
		Phone:     r.Phone,
		Address:   r.Address,
		AvatarUrl: r.AvatarURL,
		IsDeleted: r.IsDeleted,
	}
}
