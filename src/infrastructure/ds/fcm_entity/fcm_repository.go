package fcm_entity

import (
	"context"
	"github.com/hayashiki/audiy-api/src/infrastructure/ds"
	"github.com/pkg/errors"
	"go.mercari.io/datastore/boom"
)

type repository struct {
	ds.Client
}

func (r *repository) GetAll(ctx context.Context) ([]*fcm, error) {
	var entities []*fcm

	if err := r.Client.GetAll(ctx, kind, nil, entities, "", 100, "created_at"); err != nil {
		return nil, errors.WithStack(err)
	}
	return entities, nil
}

func (r *repository) Get(ctx context.Context, id string) (*fcm, error) {
	entity := onlyID(id)

	if err := r.Client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity, nil
}

func (r *repository) Put(tx *boom.Transaction, item *fcm) error {
	if _, err := tx.Put(item); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
