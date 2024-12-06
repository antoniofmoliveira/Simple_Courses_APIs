package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/antoniofmoliveira/courses/flatbuffersapi/fb"
	"golang.org/x/exp/slog"
)

func main() {
	fmt.Println("### list categories")
	listCategories("http://localhost:8088/categories")
	fmt.Println("/n### list categories error")
	listCategories("http://localhost:8088/categorieserror")
	fmt.Println("/n### list categories panic")
	listCategories("http://localhost:8088/categoriespanic")

}

func listCategories(url string) {
	defer func() {
		if r := recover(); r != nil {
			slog.Info("server response doesn't match expected")
			slog.Error("listCategoriesError", "msg", r)
		}
	}()
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		slog.Error("could not create request: ", err)
	}

	req.Header.Set("Accept", "application/octet-stream")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("could not send request: ", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("could not read response: ", err)
	}

	if res.StatusCode != http.StatusOK {
		slog.Error("listCategoriesError", "response status code", res.StatusCode)
		fbMessage := fb.GetRootAsMessage(body, 0)
		slog.Info("listCategoriesError", "message", fbMessage.Message())
		slog.Info("listCategoriesError", "isSuccess", fbMessage.IsSuccess())
		return
	}

	fbCategoriesOutput := fb.GetRootAsCategories(body, 0)
	slog.Info("# of categories: %d", fbCategoriesOutput.ElementsLength())

	for i := 0; i < fbCategoriesOutput.ElementsLength(); i++ {
		var c fb.Category
		fbCategoriesOutput.Elements(&c, i)
		fmt.Printf("Category %d: id: %s, name: %s, description: %s\n", i+1, c.Id(), c.Name(), c.Description())
	}
}
