package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/ZQCard/kbk-administrator/api/administrator/v1"
	"github.com/ZQCard/kbk-administrator/internal/biz"
	"github.com/ZQCard/kbk-administrator/internal/domain"
)

type AdministratorService struct {
	v1.UnimplementedAdministratorServiceServer
	administratorUsecase *biz.AdministratorUsecase
	log                  *log.Helper
}

func NewAdministratorService(administratorUsecase *biz.AdministratorUsecase, logger log.Logger) *AdministratorService {
	return &AdministratorService{
		log:                  log.NewHelper(log.With(logger, "module", "service/AdministratorService")),
		administratorUsecase: administratorUsecase,
	}
}

func (s *AdministratorService) GetAdministratorList(ctx context.Context, reqData *v1.GetAdministratorListReq) (*v1.GetAdministratorListPageRes, error) {
	params := make(map[string]interface{})
	params["username"] = reqData.Username
	params["mobile"] = reqData.Mobile
	params["nickname"] = reqData.Nickname
	if reqData.Status != nil {
		params["status"] = reqData.Status.Value
	}
	params["created_at_start"] = reqData.CreatedAtStart
	params["created_at_end"] = reqData.CreatedAtEnd
	list, count, err := s.administratorUsecase.ListAdministrator(ctx, reqData.Page, reqData.PageSize, params)
	if err != nil {
		return nil, err
	}
	res := &v1.GetAdministratorListPageRes{}
	res.Total = int64(count)
	for _, v := range list {
		res.List = append(res.List, toDomainAdministrator(v))
	}
	return res, nil
}

func (s *AdministratorService) GetAdministrator(ctx context.Context, req *v1.GetAdministratorReq) (*v1.Administrator, error) {
	params := map[string]interface{}{}
	params["id"] = req.Id
	params["username"] = req.Username
	params["mobile"] = req.Mobile
	params["role"] = req.Role
	res, err := s.administratorUsecase.GetAdministrator(ctx, params)
	if err != nil {
		return nil, err
	}
	return toDomainAdministrator(res), nil
}

func (s *AdministratorService) CreateAdministrator(ctx context.Context, req *v1.CreateAdministratorReq) (*v1.Administrator, error) {
	res, err := s.administratorUsecase.CreateAdministrator(ctx, &domain.Administrator{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: req.Password,
		Mobile:   req.Mobile,
		Status:   req.Status,
		Avatar:   req.Avatar,
		Role:     req.Role,
	})
	if err != nil {
		return nil, err
	}
	return toDomainAdministrator(res), nil
}

func (s *AdministratorService) UpdateAdministrator(ctx context.Context, req *v1.UpdateAdministratorReq) (*emptypb.Empty, error) {
	err := s.administratorUsecase.UpdateAdministrator(ctx, &domain.Administrator{
		Id:       req.Id,
		Username: req.Username,
		Password: req.Password,
		Mobile:   req.Mobile,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Role:     req.Role,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AdministratorService) DeleteAdministrator(ctx context.Context, req *v1.DeleteAdministratorReq) (*emptypb.Empty, error) {
	err := s.administratorUsecase.DeleteAdministrator(ctx, &domain.Administrator{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AdministratorService) RecoverAdministrator(ctx context.Context, req *v1.RecoverAdministratorReq) (*emptypb.Empty, error) {
	err := s.administratorUsecase.RecoverAdministrator(ctx, &domain.Administrator{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AdministratorService) VerifyAdministratorPassword(ctx context.Context, req *v1.VerifyAdministratorPasswordReq) (*emptypb.Empty, error) {
	ok, err := s.administratorUsecase.VerifyAdministratorPassword(ctx, req.Id, req.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &emptypb.Empty{}, v1.ErrorBadRequest("密码错误")
	}
	return &emptypb.Empty{}, nil
}

func (s *AdministratorService) AdministratorStatusChange(ctx context.Context, req *v1.AdministratorStatusChangeReq) (*emptypb.Empty, error) {
	ok, err := s.administratorUsecase.AdministratorStatusChange(ctx, req.Id, req.Status)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &emptypb.Empty{}, v1.ErrorSystemError("状态更改失败")
	}
	return &emptypb.Empty{}, nil
}

func (s *AdministratorService) AdministratorLoginSuccess(ctx context.Context, req *v1.AdministratorLoginSuccessReq) (*emptypb.Empty, error) {
	ok, err := s.administratorUsecase.UpdateAdministratorLoginInfo(ctx, req.Id, req.LastLoginIp, req.LastLoginTime)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &emptypb.Empty{}, v1.ErrorSystemError("更新登录信息失败")
	}
	return &emptypb.Empty{}, nil
}

func toDomainAdministrator(administrator *domain.Administrator) *v1.Administrator {
	return &v1.Administrator{
		Id:        administrator.Id,
		Username:  administrator.Username,
		Mobile:    administrator.Mobile,
		Nickname:  administrator.Nickname,
		Avatar:    administrator.Avatar,
		Status:    administrator.Status,
		Role:      administrator.Role,
		CreatedAt: administrator.CreatedAt,
		UpdatedAt: administrator.UpdatedAt,
	}
}
