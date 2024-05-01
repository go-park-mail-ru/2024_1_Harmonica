package repository

import (
	"context"
	"harmonica/internal/entity"
)

func (r *DBRepository) SearchForUsers(ctx context.Context, query string) ([]entity.SearchUser, error) {
	return []entity.SearchUser{}, nil
}

func (r *DBRepository) SearchForPins(ctx context.Context, query string) ([]entity.SearchPin, error) {
	return []entity.SearchPin{}, nil
}

func (r *DBRepository) SearchForBoards(ctx context.Context, query string) ([]entity.SearchBoard, error) {
	return []entity.SearchBoard{}, nil
}
