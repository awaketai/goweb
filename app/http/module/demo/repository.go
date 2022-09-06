package demo

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetUserIds() []int {
	return []int{12, 4}
}

func (r *Repository) GetUserByIds([]int) []UserModel {
	return []UserModel{
		{
			UserId: 12,
			Name:   "foo",
			Age:    11,
		},
		{
			UserId: 4,
			Name:   "bar",
			Age:    15,
		},
	}

}
