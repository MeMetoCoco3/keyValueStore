package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type BoltServer[K comparable, V any] struct {
	//*Server
	Storage    BoltStorer[K, V]
	ListenAddr string
}

// Constructors cannot work with generics
func NewBoltServer(boltPath, lAddr string) (*BoltServer[int, *User], error) {
	boltStore, err := NewBoltStore[int, *User](boltPath, "Bunny")
	if err != nil {
		return nil, fmt.Errorf("Error creating new bolt server. %v", err)
	}

	return &BoltServer[int, *User]{
		Storage:    boltStore,
		ListenAddr: lAddr,
	}, nil
}

func (s *BoltServer[K, V]) StartBoltEcho() error {
	fmt.Printf("We starting the echo serer on port %s", s.ListenAddr)
	e := echo.New()
	e.GET("/put/:id/:name/:age/:role", s.HandlePutEcho)
	e.GET("/get/:id", s.HandleGetEcho)
	e.GET("/getAll", s.HandleGetAllEcho)
	e.DELETE("/delete/:id", s.HandleDeleteEcho)
	return e.Start(s.ListenAddr)
}
