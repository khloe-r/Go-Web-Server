package main

import (
	"fmt"
  "log"
  "net/http"
  "html/template"
)

type ToDoList struct {
  toDoItems []string
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (todolist *ToDoList) formHandler(w http.ResponseWriter, r *http.Request) {
  if err := r.ParseForm(); err != nil {
    fmt.Fprintf(w, "ParseForm() err: %v", err)
    return
  }
  
  item := r.FormValue("item")
  todolist.toDoItems = append(todolist.toDoItems, item)
  renderTemplate(w, "./static/index.html", todolist.toDoItems)
}

func (todolist *ToDoList) deleteHandler(w http.ResponseWriter, r *http.Request) {
  if err := r.ParseForm(); err != nil {
    fmt.Fprintf(w, "ParseForm() err: %v", err)
    return
  }

  item := r.FormValue("item")
  for i, other := range todolist.toDoItems {
      if other == item {
          todolist.toDoItems = append(todolist.toDoItems[:i], todolist.toDoItems[i+1:]...)
          break
      }
  }
  renderTemplate(w, "./static/index.html", todolist.toDoItems)
}

func main() {
  items := []string{}
  list := ToDoList{toDoItems: items}

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "./static/index.html", list.toDoItems)
	})
  http.HandleFunc("/add-todo", list.formHandler)
  http.HandleFunc("/delete-todo", list.deleteHandler)

  fmt.Printf("Starting server at port 8080\n")
  
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal(err)
  }
}
