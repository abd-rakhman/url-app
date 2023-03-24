package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/abd-rakhman/url-app/utils"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, txFunc func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = txFunc(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error in execTx %v, error in rollback %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type CreateNewURLRequest struct {
	URL    string
	HashID string
}

func (store *Store) CreateNewURLTx(ctx context.Context, args CreateNewURLRequest) (Url, error) {
	var result Url
	err := store.execTx(ctx, func(q *Queries) error {
		if args.HashID != "" {
			_, err := q.GetUrlByHashIdForUpdate(ctx, args.HashID)
			if err != nil {
				if err == sql.ErrNoRows {
					result, err = q.CreateUrl(ctx, CreateUrlParams{
						HashID: args.HashID,
						Url:    args.URL,
					})
					if err != nil {
						return err
					} else {
						return nil
					}
				} else {
					return err
				}
			} else {
				return fmt.Errorf("hash id already exists")
			}
		} else {
			for hashedLength := 6; hashedLength < 10; hashedLength++ {
				for repeat := 0; repeat < 4; repeat++ {
					hashId := utils.RandomString(hashedLength)
					_, err := q.GetUrlByHashIdForUpdate(ctx, hashId)
					if err != nil {
						if err == sql.ErrNoRows {
							result, err = q.CreateUrl(ctx, CreateUrlParams{
								HashID: hashId,
								Url:    args.URL,
							})
							if err != nil {
								return err
							} else {
								return nil
							}
						} else {
							return err
						}
					}
				}
			}
		}
		return fmt.Errorf("internal server error: unable to create new url")
	})
	return result, err
}
