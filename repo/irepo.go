package repo

import (
	"context"
	"github.com/ervitis/spamtoputocorreos"
)

type (
	IRepository interface {
		Save(context.Context, *spamtoputocorreos.StatusTrace) error
		Get(context.Context, string) (*spamtoputocorreos.StatusTrace, error)
		Delete(context.Context) error
	}
)
