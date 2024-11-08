package resolver_user

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	graph_model "e-learning/src/graph/generated/model"
	generated_user "e-learning/src/graph/generated/user"
	"fmt"
)

// FindAccountByID is the resolver for the findAccountByID field.
func (r *entityResolver) FindAccountByID(ctx context.Context, id string) (*graph_model.Account, error) {
	panic(fmt.Errorf("not implemented: FindAccountByID - findAccountByID"))
}

// FindUserInforByID is the resolver for the findUserInforByID field.
func (r *entityResolver) FindUserInforByID(ctx context.Context, id string) (*graph_model.UserInfor, error) {
	panic(fmt.Errorf("not implemented: FindUserInforByID - findUserInforByID"))
}

// Entity returns generated_user.EntityResolver implementation.
func (r *Resolver) Entity() generated_user.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
