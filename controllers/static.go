package controllers

import (
	"html/template"
	"net/http"
)

type Static struct {
	Template Templater
}

func StaticHandler(tpl Templater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Templater) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! We offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails 24/7, though response",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="mailto:support@asd.com">support@asd.com</a>`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
