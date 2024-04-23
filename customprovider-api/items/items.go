package items

import (
	"errors"

	"github.com/google/uuid"
)

type ItemsResponse struct {
	Items []*Item `json:"items"`
}

type Items struct {
	Items map[string]*Item `json:"items"`
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func New() *Items {
	return &Items{
		Items: make(map[string]*Item),
	}
}

// Create
func (s *Items) AddItem(name string) *Item {
	newID := uuid.New().String()
	s.Items[newID] = &Item{
		ID:   newID,
		Name: name,
	}
	return s.Items[newID]
}

// Read
func (s *Items) ReadItems() ItemsResponse {
	var itemsResp ItemsResponse
	for _, item := range s.Items {
		itemsResp.Items = append(itemsResp.Items, item)
	}
	return itemsResp
}

func (s *Items) ReadItemByID(id string) (*Item, error) {
	item, ok := s.Items[id]
	if !ok {
		return nil, errors.New("No item with that id")
	}
	return item, nil
}

// Update
func (s *Items) Update(id string, updatedName string) (*Item, error) {
	_, ok := s.Items[id]
	if !ok {
		return nil, errors.New("No item with that id")
	}
	s.Items[id] = &Item{
		ID:   id,
		Name: updatedName,
	}
	return s.Items[id], nil
}

// Delete
func (s *Items) Remove(id string) {
	delete(s.Items, id)
}
