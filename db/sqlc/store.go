package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	URL       string
	HashID    string
	ExpiresAt int64
}

func (store *Store) CreateNewURLTx(ctx context.Context, args CreateNewURLRequest) (Url, error) {
	var result Url
	err := store.execTx(ctx, func(q *Queries) error {
		if args.HashID != "" {
			_, err := q.GetUrlByHashIdForUpdate(ctx, args.HashID)
			if err != nil {
				if err == sql.ErrNoRows {
					result, err = q.createURL(ctx, args)
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
							result, err = q.createURL(ctx, CreateNewURLRequest{
								URL:       args.URL,
								HashID:    hashId,
								ExpiresAt: args.ExpiresAt,
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

func (q *Queries) createURL(ctx context.Context, args CreateNewURLRequest) (Url, error) {
	if args.ExpiresAt == 0 {
		args := CreateUrlParams{
			HashID: args.HashID,
			Url:    args.URL,
		}
		return q.CreateUrl(ctx, args)
	} else {
		args := CreateUrlWithExpiresAtParams{
			HashID:    args.HashID,
			Url:       args.URL,
			ExpiresAt: time.Unix(args.ExpiresAt, 0).UTC(),
		}
		fmt.Printf("Create with args %v\n", args)
		return q.CreateUrlWithExpiresAt(ctx, args)
	}
}
