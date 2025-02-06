package gorm

import (
	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository/dao/entity"
	"context"
	"database/sql"
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		sqlmock func(t *testing.T) *sql.DB

		ctx  context.Context
		user entity.User

		wantErr error
	}{
		{
			name: "insert success",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mockRes := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			user: entity.User{
				Email: sql.NullString{
					String: "test@test.com",
					Valid:  true,
				},
				Password: "test1234",
			},

			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
					})
				return db
			},
			ctx:  context.Background(),
			user: entity.User{},

			wantErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(errors.New("数据库错误"))
				return db
			},
			ctx:  context.Background(),
			user: entity.User{},

			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sqlDB := tc.sqlmock(t)
			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewUserDao(db)
			err = dao.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
