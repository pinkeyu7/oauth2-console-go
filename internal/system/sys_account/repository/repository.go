package repository

import (
	"oauth2-console-go/dto/model"
	"oauth2-console-go/internal/system/sys_account"

	"xorm.io/xorm"
)

type Repository struct {
	orm *xorm.EngineGroup
}

func NewRepository(orm *xorm.EngineGroup) sys_account.Repository {
	return &Repository{
		orm: orm,
	}
}

func (r *Repository) Insert(m *model.SysAccount) error {
	_, err := r.orm.Insert(m)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Find(offset, limit int) ([]*model.SysAccount, error) {
	list := make([]*model.SysAccount, 0)

	err := r.orm.Limit(limit, offset).Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *Repository) FindOne(m *model.SysAccount) (*model.SysAccount, error) {
	has, err := r.orm.Get(m)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return m, nil
}

func (r *Repository) Count() (int, error) {
	count, err := r.orm.Count(&model.SysAccount{})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *Repository) Update(m *model.SysAccount) error {
	_, err := r.orm.ID(m.Id).Update(m)
	if err != nil {
		return err
	}
	return nil
}
