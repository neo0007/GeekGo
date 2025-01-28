package repository

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/domain"
	"context"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUserRepositoryDao_FindById(t *testing.T) {
	testCases := []struct {
		name string
		ctx  context.Context
		id   int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "success",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

		})
	}
}
