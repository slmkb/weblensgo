package main

import (
	"html/template"
	"os"
)

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	u := struct {
		Name   string
		Age    int
		Nested struct {
			NestedMap   map[string]int
			NestedSlice []string
		}
	}{
		Name: "Kabekaes",
		Age:  9001,
		Nested: struct {
			NestedMap   map[string]int
			NestedSlice []string
		}{
			NestedMap: map[string]int{
				"Value1": 33,
				"Value2": 44,
			},
			NestedSlice: []string{
				"SliceString1",
				"SliceString2",
			},
		},
	}

	if err = t.Execute(os.Stdout, u); err != nil {
		panic(err)
	}
}
