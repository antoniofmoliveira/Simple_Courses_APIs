package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/antoniofmoliveira/courses/flatbuffersapi/fb"
)

func main() {

	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:8088/categories", nil)
	if err != nil {
		log.Println("could not create request: ", err)
	}

	req.Header.Set("Accept", "application/octet-stream")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("could not send request: ", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("could not read response: ", err)
	}

	fbCategoriesOutput := fb.GetRootAsCategories(body, 0)
	log.Printf("categories: %d", fbCategoriesOutput.ElementsLength())
	// for i := 0; i < fbCategoriesOutput.ElementsLength(); i++ {
	category := new(fb.Category)
	if fbCategoriesOutput.Elements(category, 0) {
		id := string(category.Id())
		name := string(category.Name())
		description := string(category.Description())
		log.Printf("category: id: %s, name: %s, description: %s", id, name, description)
	}
	// }

}
