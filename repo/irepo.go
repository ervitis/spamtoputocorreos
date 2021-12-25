package repo

import (
	"context"
	"github.com/ervitis/spamtoputocorreos/models"
)

type (
	IRepository interface {
		Save(context.Context, *models.StatusTrace) error
		Get(context.Context, string) (*models.StatusTrace, error)
		Delete(context.Context) error
	}
)
