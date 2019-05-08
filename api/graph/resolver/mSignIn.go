package resolver

import (
	"context"

	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// SignIn resolves Mutation.signIn
func (rsv *Resolver) SignIn(
	ctx context.Context,
	params struct {
		Email    string
		Password string
	},
) (*Session, error) {
	// Validate inputs
	if len(params.Email) < 1 || len(params.Password) < 1 {
		err := strerr.New(strerr.ErrInvalidInput, "missing credentials")
		rsv.error(ctx, err)
		return nil, err
	}

	transactRes, err := rsv.str.CreateSession(
		ctx,
		params.Email,
		params.Password,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil, err
	}

	return &Session{
		root:     rsv,
		uid:      transactRes.UID,
		key:      transactRes.Key,
		creation: transactRes.CreationTime,
		userUID:  transactRes.UserUID,
	}, nil
}
