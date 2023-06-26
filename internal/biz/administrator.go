package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	v1 "github.com/ZQCard/kratos-base-kit/kbk-administrator/api/administrator/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-administrator/internal/domain"
	"github.com/ZQCard/kratos-base-kit/kbk-administrator/pkg/utils/encryptHelper"
)

type AdministratorRepo interface {
	ListAdministrator(ctx context.Context, page, pageSize int64, params map[string]interface{}) ([]*domain.Administrator, int64, error)
	CreateAdministrator(ctx context.Context, administrator *domain.Administrator) (*domain.Administrator, error)
	UpdateAdministrator(ctx context.Context, administrator *domain.Administrator) error
	DeleteAdministrator(ctx context.Context, administrator *domain.Administrator) error
	RecoverAdministrator(ctx context.Context, administrator *domain.Administrator) error
	GetAdministrator(ctx context.Context, params map[string]interface{}) (*domain.Administrator, error)
	VerifyAdministratorPassword(ctx context.Context, id int64, password string) (bool, error)
	AdministratorStatusChange(ctx context.Context, id int64, status bool) (bool, error)
	UpdateAdministratorLoginInfo(ctx context.Context, id int64, ip string, time string) (bool, error)
}

type AdministratorUsecase struct {
	repo   AdministratorRepo
	logger *log.Helper
}

func NewAdministratorUsecase(repo AdministratorRepo, logger log.Logger) *AdministratorUsecase {
	return &AdministratorUsecase{repo: repo, logger: log.NewHelper(log.With(logger, "module", "usecase/administrator"))}
}

func (suc *AdministratorUsecase) ListAdministrator(ctx context.Context, page, pageSize int64, params map[string]interface{}) ([]*domain.Administrator, int64, error) {
	return suc.repo.ListAdministrator(ctx, page, pageSize, params)
}

func (suc *AdministratorUsecase) CreateAdministrator(ctx context.Context, administrator *domain.Administrator) (*domain.Administrator, error) {
	// 查看用户名是否存在
	recordTmp, _ := suc.repo.GetAdministrator(ctx, map[string]interface{}{
		"username": administrator.Username,
	})

	if recordTmp != nil {
		return nil, v1.ErrorRecordAlreadyExists("用户名已存在")
	}
	salt, password := encryptHelper.HashPassword(administrator.Password)
	administrator.Salt = salt
	administrator.Password = password
	return suc.repo.CreateAdministrator(ctx, administrator)
}

func (suc *AdministratorUsecase) GetAdministrator(ctx context.Context, params map[string]interface{}) (*domain.Administrator, error) {
	return suc.repo.GetAdministrator(ctx, params)
}

func (suc *AdministratorUsecase) UpdateAdministrator(ctx context.Context, administrator *domain.Administrator) error {
	recordTmp, err := suc.repo.GetAdministrator(ctx, map[string]interface{}{
		"id": administrator.Id,
	})
	if err != nil {
		return err
	}
	// 如果有密码更改
	if administrator.Password != "" {
		salt, password := encryptHelper.HashPassword(administrator.Password)
		administrator.Salt = salt
		administrator.Password = password
	} else {
		administrator.Salt = recordTmp.Salt
		administrator.Password = recordTmp.Password
	}
	return suc.repo.UpdateAdministrator(ctx, administrator)
}

func (suc *AdministratorUsecase) DeleteAdministrator(ctx context.Context, administrator *domain.Administrator) error {
	return suc.repo.DeleteAdministrator(ctx, administrator)
}

func (suc *AdministratorUsecase) RecoverAdministrator(ctx context.Context, administrator *domain.Administrator) error {
	return suc.repo.RecoverAdministrator(ctx, administrator)
}

func (suc *AdministratorUsecase) VerifyAdministratorPassword(ctx context.Context, id int64, password string) (bool, error) {
	return suc.repo.VerifyAdministratorPassword(ctx, id, password)
}

func (suc *AdministratorUsecase) AdministratorStatusChange(ctx context.Context, id int64, status bool) (bool, error) {
	_, err := suc.repo.GetAdministrator(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return false, err
	}
	return suc.repo.AdministratorStatusChange(ctx, id, status)
}

func (suc *AdministratorUsecase) UpdateAdministratorLoginInfo(ctx context.Context, id int64, ip string, time string) (bool, error) {
	return suc.repo.UpdateAdministratorLoginInfo(ctx, id, ip, time)
}
