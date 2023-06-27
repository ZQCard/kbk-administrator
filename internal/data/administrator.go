package data

import (
	"context"

	v1 "github.com/ZQCard/kbk-administrator/api/administrator/v1"
	"github.com/ZQCard/kbk-administrator/internal/biz"
	"github.com/ZQCard/kbk-administrator/internal/domain"
	"github.com/ZQCard/kbk-administrator/pkg/utils/encryptHelper"
	"github.com/ZQCard/kbk-administrator/pkg/utils/timeHelper"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type AdministratorEntity struct {
	BaseFields
	Domain        string `gorm:"type:varchar(255);not null;comment:域"`
	Status        bool   `gorm:"not null;comment:状态0冻结1正常"`
	Username      string `gorm:"type:varchar(255);not null;unique;comment:用户名"`
	Password      string `gorm:"type:varchar(255);not null;comment:密码"`
	Salt          string `gorm:"type:varchar(255);not null;comment:密码盐"`
	Mobile        string `gorm:"type:varchar(255);not null;comment:手机号"`
	Nickname      string `gorm:"type:varchar(255);not null;comment:昵称"`
	Avatar        string `gorm:"type:varchar(255);not null;comment:头像"`
	Role          string `gorm:"type:varchar(255);not null;comment:角色"`
	LastLoginTime string `gorm:"type:varchar(255);not null;comment:上次登录时间"`
	LastLoginIp   string `gorm:"type:varchar(255);not null;comment:上次登录ip"`
}

func (AdministratorEntity) TableName() string {
	return "administrators"
}

type AdministratorRepo struct {
	data *Data
	log  *log.Helper
}

func NewAdministratorRepo(data *Data, logger log.Logger) biz.AdministratorRepo {
	repo := &AdministratorRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/administrator")),
	}
	return repo
}

// searchParam 搜索条件
func (repo AdministratorRepo) searchParam(ctx context.Context, params map[string]interface{}) *gorm.DB {
	conn := repo.data.db.Model(&AdministratorEntity{})
	if id, ok := params["id"]; ok && id != nil && id.(int64) != 0 {
		conn = conn.Where("id = ?", id)
	}
	if mobile, ok := params["mobile"]; ok && mobile.(string) != "" {
		conn = conn.Where("mobile LIKE ?", "%"+mobile.(string)+"%")
	}
	// 开始时间
	if start, ok := params["created_at_start"]; ok && start.(string) != "" {
		conn = conn.Where("created_at >= ?", start.(string)+" 00:00:00")
	}
	// 结束时间
	if end, ok := params["created_at_end"]; ok && end.(string) != "" {
		conn = conn.Where("created_at <= ?", end.(string)+" 23:59:59")
	}
	if nickname, ok := params["nickname_like"]; ok && nickname != nil && nickname.(string) != "" {
		conn = conn.Where("nickname LIKE ?", "%"+nickname.(string)+"%")
	}
	if nickname, ok := params["nickname"]; ok && nickname != nil && nickname.(string) != "" {
		conn = conn.Where("nickname = ?", nickname)
	}
	if username, ok := params["username_like"]; ok && username != nil && username.(string) != "" {
		conn = conn.Where("username LIKE ?", "%"+username.(string)+"%")
	}
	if username, ok := params["username"]; ok && username != nil && username.(string) != "" {
		conn = conn.Where("username = ?", username)
	}
	if status, ok := params["status"]; ok && status != nil {
		conn = conn.Where("status = ?", status)
	}
	if role, ok := params["role"]; ok && role != nil && role.(string) != "" {
		conn = conn.Where("role = ?", role)
	}
	getDbWithDomain(ctx, conn)
	return conn
}

func (repo AdministratorRepo) GetAdministratorByParams(ctx context.Context, params map[string]interface{}) (record *AdministratorEntity, err error) {
	if len(params) == 0 {
		return &AdministratorEntity{}, v1.ErrorBadRequest("缺少搜索条件")
	}
	conn := repo.searchParam(ctx, params)
	if err = conn.First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &AdministratorEntity{}, v1.ErrorRecordNotFound("数据不存在")
		}
		return record, v1.ErrorSystemError("GetAdministratorByParams First Error : %s", err)
	}
	return record, nil
}

func (repo AdministratorRepo) CreateAdministrator(ctx context.Context, domain *domain.Administrator) (*domain.Administrator, error) {
	entity := &AdministratorEntity{
		Domain:   getDomain(ctx),
		Username: domain.Username,
		Salt:     domain.Salt,
		Password: domain.Password,
		Nickname: domain.Nickname,
		Mobile:   domain.Mobile,
		Status:   domain.Status,
		Role:     domain.Role,
		Avatar:   domain.Avatar,
	}
	if err := repo.data.db.Model(entity).Create(entity).Error; err != nil {
		return nil, v1.ErrorSystemError("CreateAdministrator Create Error : %s", err)
	}
	response := toDomainAdministrator(entity)
	return response, nil
}

func (repo AdministratorRepo) UpdateAdministrator(ctx context.Context, domain *domain.Administrator) error {
	// 根据Id查找记录
	record, err := repo.GetAdministratorByParams(ctx, map[string]interface{}{
		"id": domain.Id,
	})
	if err != nil {
		return err
	}
	record.Domain = getDomain(ctx)
	record.Salt = domain.Salt
	record.Password = domain.Password
	record.Avatar = domain.Avatar
	record.Nickname = domain.Nickname
	record.Status = domain.Status
	record.Mobile = domain.Mobile
	record.Role = domain.Role
	if err := getDbWithDomain(ctx, repo.data.db).Model(&record).Where("id = ?", record.Id).Save(&record).Error; err != nil {
		return v1.ErrorSystemError("UpdateAdministrator Save Error : %s", err.Error())
	}
	return nil
}

func (repo AdministratorRepo) GetAdministrator(ctx context.Context, params map[string]interface{}) (*domain.Administrator, error) {
	record, err := repo.GetAdministratorByParams(ctx, params)
	if err != nil {
		return nil, err
	}
	// 返回数据
	response := toDomainAdministrator(record)
	return response, nil
}

func (repo AdministratorRepo) ListAdministrator(ctx context.Context, page, pageSize int64, params map[string]interface{}) ([]*domain.Administrator, int64, error) {
	list := []*AdministratorEntity{}
	conn := repo.searchParam(ctx, params)
	err := conn.Scopes(Paginate(page, pageSize)).Find(&list).Error
	if err != nil {
		return nil, 0, v1.ErrorSystemError("ListAdministrator Find Error : %s", err.Error())
	}

	count := int64(0)
	conn.Count(&count)
	rv := make([]*domain.Administrator, 0, len(list))
	for _, record := range list {
		administrator := toDomainAdministrator(record)
		rv = append(rv, administrator)
	}
	return rv, count, nil
}

func (repo AdministratorRepo) DeleteAdministrator(ctx context.Context, domain *domain.Administrator) error {
	err := getDbWithDomain(ctx, repo.data.db).Where("id = ?", domain.Id).Delete(&AdministratorEntity{}).Error
	if err != nil {
		return v1.ErrorSystemError("DeleteAdministrator Delete Error : %s", err)
	}
	return nil
}

func (repo AdministratorRepo) RecoverAdministrator(ctx context.Context, domain *domain.Administrator) error {
	err := getDbWithDomain(ctx, repo.data.db).Unscoped().Model(AdministratorEntity{}).Where("id = ?", domain.Id).UpdateColumn("deleted_at", nil).Error
	if err != nil {
		return v1.ErrorSystemError("RecoverAdministrator UpdateColumn Error : %s", err)
	}
	return nil
}

func (repo AdministratorRepo) VerifyAdministratorPassword(ctx context.Context, id int64, password string) (bool, error) {
	recordTmp, err := repo.GetAdministrator(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return false, err
	}
	// 验证密码
	if !encryptHelper.CheckPasswordHash(password, recordTmp.Salt, recordTmp.Password) {
		return false, nil
	}
	return true, nil
}

func (repo AdministratorRepo) AdministratorStatusChange(ctx context.Context, id int64, status bool) (bool, error) {
	success := getDbWithDomain(ctx, repo.data.db).Model(AdministratorEntity{}).Where("id = ?", id).UpdateColumn("status", status).RowsAffected
	if success == 0 {
		return false, nil
	}
	return true, nil
}

func (repo AdministratorRepo) UpdateAdministratorLoginInfo(ctx context.Context, id int64, ip string, time string) (bool, error) {
	success := getDbWithDomain(ctx, repo.data.db).Model(AdministratorEntity{}).Where("id = ?", id).UpdateColumns(
		map[string]interface{}{
			"last_login_ip":   ip,
			"last_login_time": time,
		},
	).RowsAffected
	if success == 0 {
		return false, nil
	}
	return true, nil
}

func toDomainAdministrator(administrator *AdministratorEntity) *domain.Administrator {
	return &domain.Administrator{
		Id:        administrator.Id,
		Salt:      administrator.Salt,
		Password:  administrator.Password,
		Username:  administrator.Username,
		Mobile:    administrator.Mobile,
		Nickname:  administrator.Nickname,
		Avatar:    administrator.Avatar,
		Status:    administrator.Status,
		Role:      administrator.Role,
		CreatedAt: timeHelper.FormatYMDHIS(&administrator.CreatedAt),
		UpdatedAt: timeHelper.FormatYMDHIS(&administrator.UpdatedAt),
	}
}
