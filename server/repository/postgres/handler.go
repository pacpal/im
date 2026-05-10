// 后续添加其他数据库
package postgres

import (
	"IM/server/model"
	repo "IM/server/repository"
	"context"

	"gorm.io/gorm"
)

var _ repo.UserRepo = (*UserRepoImpl)(nil)
var _ repo.GroupRepo = (*GroupRepoImpl)(nil)
var _ repo.MsgRepo = (*MsgRepoImpl)(nil)

type UserRepoImpl struct {
	db *gorm.DB
}
type GroupRepoImpl struct {
	db *gorm.DB
}
type MsgRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repo.UserRepo {
	return &UserRepoImpl{db: db}
}
func GroupUserRepo(db *gorm.DB) repo.UserRepo {
	return &GroupRepoImpl{db: db}
}
func NewMsgRepo(db *gorm.DB) repo.UserRepo {
	return &UserMsgImpl{db: db}
}
func (r *UserRepoImpl) GetUserByID(ctx context.Context, uid string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id=?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepoImpl) GetUserByTele(ctx context.Context, tele string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("tele=?", tele).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepoImpl) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("name=?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepoImpl) Create(ctx context.Context, uid, name string, password []byte) (bool, error) {
	user := &model.User{
		ID:       uid,
		Name:     name,
		Password: string(password),
	}
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
func (r *GroupRepoImpl)GetGroupByID(ctx context.Context, gid string) (*model.Group, error){
	var group *model.Group
	err:=r.db.WithContext(ctx).Where("id=?", gid).First(&group).Error
	if err != nil {
		return nil, err
	}
	return group, nil
}
func (r *GroupRepoImpl)IsMember(ctx context.Context, gid, uid string) (bool, error){
	var exist bool
	err:=r.db.WithContext(ctx).Where("group_id", args ...interface{})
}
func(r *GroupRepoImpl)SaveGroup(ctx context.Context, user *model.Group) error{

}

func(r *MsgRepoImpl)GetOfflineMsgs(ctx context.Context, uid string) (*[]model.Message, error){

}
func(r *MsgRepoImpl)ClearOfflineMsgs(ctx context.Context, uid string){

}
func(r *MsgRepoImpl)GetGroupByID(ctx context.Context, gid string) (*model.Group, error){

}