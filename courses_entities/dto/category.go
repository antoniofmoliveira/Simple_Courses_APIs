package dto

type CategoryInputDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryOutputDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryListOutputDto struct {
	Categories []CategoryOutputDto `json:"categories"`
}
