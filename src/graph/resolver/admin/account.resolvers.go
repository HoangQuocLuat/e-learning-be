package resolver_admin

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	generated_admin "e-learning/src/graph/generated/admin"
	graph_model "e-learning/src/graph/generated/model"
	service_account "e-learning/src/service/account"
	"fmt"
)

// AccountAdd is the resolver for the accountAdd field.
func (r *mutationResolver) AccountAdd(ctx context.Context, data *graph_model.AccountAdd) (*graph_model.Account, error) {
	input := &service_account.AccountAddCommand{
		UserName: data.Username,
		Password: data.Password,
		Role:     data.Role,
	}

	result, err := service_account.AccountAdd(ctx, input)
	if err != nil {
		return &graph_model.Account{}, err
	}

	return result.ConvertToModelGraph(), nil
}

// AccountDelete is the resolver for the AccountDelete field.
func (r *mutationResolver) AccountDelete(ctx context.Context, data *graph_model.AccountDelete) (*graph_model.Account, error) {
	panic(fmt.Errorf("not implemented: AccountDelete - AccountDelete"))
}

// AccountMe is the resolver for the accountMe field.
func (r *queryResolver) AccountMe(ctx context.Context) (*graph_model.Account, error) {
	panic(fmt.Errorf("not implemented: AccountMe - accountMe"))
}

// AccountPagination is the resolver for the accountPagination field.
func (r *queryResolver) AccountPagination(ctx context.Context, page int, limit int, orderBy *string, search map[string]interface{}) (*graph_model.AccountPagination, error) {
	panic(fmt.Errorf("not implemented: AccountPagination - accountPagination"))
}

// Mutation returns generated_admin.MutationResolver implementation.
func (r *Resolver) Mutation() generated_admin.MutationResolver { return &mutationResolver{r} }

// Query returns generated_admin.QueryResolver implementation.
func (r *Resolver) Query() generated_admin.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }