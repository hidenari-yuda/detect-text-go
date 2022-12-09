package responses

import "github.com/hidenari-yuda/paychan/domain/entity"

type User struct {
	User *entity.User `json:"user"`
}

func NewUser(user *entity.User) User {
	return User{
		User: user,
	}
}

type UserList struct {
	UserList []*entity.User `json:"user_list"`
}

func NewUserList(users []*entity.User) UserList {
	return UserList{
		UserList: users,
	}
}
