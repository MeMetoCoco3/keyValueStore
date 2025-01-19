package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) StartEcho() {
	fmt.Printf("We starting the echo serer on port %s", s.ListenAddr)
	e := echo.New()
	e.GET("/put/:id/:name/:age/:role", s.HandlePutEcho)
	e.GET("/get/:id", s.HandleGetEcho)
	e.Start(s.ListenAddr)
}

func (s *Server) HandlePutEcho(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalf("Error making conversion: %v", err)
		return err
	}

	name := c.Param("name")
	age, err := strconv.Atoi(c.Param("age"))
	if err != nil {
		log.Fatalf("Error making conversion: %v", err)
		return err
	}

	role := c.Param("role")
	newUser := NewUser(name, role, age)
	s.Storage.Put(id, newUser)

	msg := fmt.Sprintf("Added User{'name':'%s','age':%d,'role':'%s'}", name, age, role)
	return c.JSON(http.StatusOK, map[string]string{"msg": msg})
}

func (s *Server) HandleGetEcho(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalf("Error making conversion: %v", err)
		return err
	}

	user, err := s.Storage.Get(id)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}

	buff, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"msg": string(buff)})
}
