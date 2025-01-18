package main

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID
	Name string
	Role string
	Age  int
}

type UserCollection struct {
	TotalCount int
	Users      []*User
}

func NewUser(name, role string, age int) *User {
	return &User{
		ID:   uuid.New(),
		Name: name,
		Role: role,
		Age:  age,
	}
}
