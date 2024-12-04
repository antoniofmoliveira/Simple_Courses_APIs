package entity

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewCategory(id string, name, description string) *Category {
	return &Category{
		ID:          id,
		Name:        name,
		Description: description,
	}
}
