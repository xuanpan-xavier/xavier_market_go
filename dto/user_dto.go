package dto

import "xmarket_gin/model"

type UserDto struct {
	Name      string `json:"Name"`
	Telephone string `json:"Telephone"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:      user.Name,
		Telephone: user.Telephone,
	}
}
