package demo

import "github.com/awaketai/goweb/app/provider/demo"

func UserModelsToUserDTOs(models []UserModel) []UserDTO {
	ret := make([]UserDTO, len(models))
	for _, model := range models {
		t := UserDTO{
			ID:   model.UserId,
			Name: model.Name,
		}
		ret = append(ret, t)
	}
	return ret
}

func StudentToUserDTOs(students []demo.Student) []UserDTO {
	ret := make([]UserDTO, len(students))
	for _, student := range students {
		t := UserDTO{
			ID:   student.ID,
			Name: student.Name,
		}
		ret = append(ret, t)
	}
	return ret
}
