package actions

import (
  "net/http"

  "github.com/gobuffalo/buffalo"
)

type PostsResource struct{
  buffalo.Resource
}


// List default implementation.
func (v PostsResource) List(c buffalo.Context) error {
  return c.Render(http.StatusOK, r.String("Post#List"))
}


// Show default implementation.
func (v PostsResource) Show(c buffalo.Context) error {
  return c.Render(http.StatusOK, r.String("Post#Show"))
}


// Create default implementation.
func (v PostsResource) Create(c buffalo.Context) error {
  return c.Render(http.StatusOK, r.String("Post#Create"))
}


// Update default implementation.
func (v PostsResource) Update(c buffalo.Context) error {
  return c.Render(http.StatusOK, r.String("Post#Update"))
}


// Destroy default implementation.
func (v PostsResource) Destroy(c buffalo.Context) error {
  return c.Render(http.StatusOK, r.String("Post#Destroy"))
}


// New default implementation.
func (v PostsResource) New(c buffalo.Context) error {
  return c.Render(http.StatusOK, r.String("Post#New"))
}


// Edit default implementation.
func (v PostsResource) Edit(c buffalo.Context) error {
  return c.Render(http.StatusOK, r.String("Post#Edit"))
}


