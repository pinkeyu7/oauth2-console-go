package sys_account

import "oauth2-console-go/dto/model"

type Repository interface {
	Insert(m *model.SysAccount) error
	Find(offset, limit int) ([]*model.SysAccount, error)
	FindOne(m *model.SysAccount) (*model.SysAccount, error)
	Count() (int, error)
	Update(m *model.SysAccount) error
}
