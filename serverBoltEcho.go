package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

/*
	func (s *Server) StartEcho() {
		fmt.Printf("We starting the echo serer on port %s", s.ListenAddr)
		e := echo.New()
		e.GET("/put/:id/:name/:age/:role", s.HandlePutEcho)
		e.GET("/get/:id", s.HandleGetEcho)
		e.Start(s.ListenAddr)
	}
*/

// Methods Can Work with generics
/*
func (b *BoltServer[K, V]) HandlePut(w http.ResponseWriter, r *http.Request) {
	var buff bytes.Buffer

	err := IssueList.Execute(&buff, list)
	if err != nil {
		log.Println("We got error %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buff.Bytes())
}
*/

func (b *BoltServer[K, V]) HandlePutEcho(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Error converting id:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Invalid id"})
	}

	name := c.Param("name")
	age, err := strconv.Atoi(c.Param("age"))
	if err != nil {
		log.Println("Error converting age:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Invalid age"})
	}

	role := c.Param("role")
	newUser := NewUser(name, role, age)

	// Type assertion to use BoltStore methods
	if bStorage, ok := b.Storage.(BoltStorer[int, *User]); ok {
		if err := bStorage.PutB(id, newUser); err != nil {
			log.Println("Error saving user:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "Error saving user"})
		}
	} else {
		log.Println("Error casting Storage to BoltStorer")
		return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "Internal storage error"})
	}

	msg := fmt.Sprintf("Added User{'name':'%s','age':%d,'role':'%s'}", name, age, role)
	return c.JSON(http.StatusOK, map[string]string{"msg": msg})
}

func (b *BoltServer[K, V]) HandleGetEcho(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Error converting id:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Invalid id"})
	}

	var user *User
	if boltStorer, ok := b.Storage.(BoltStorer[int, *User]); ok {
		user, err = boltStorer.GetB(id)
		if err != nil {
			log.Println("Error fetching user:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "User not found"})
		}
	} else {
		log.Println("Error casting Storage to BoltStorer")
		return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "Internal storage error"})
	}

	buff, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshaling user:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "Error processing user data"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"user": json.RawMessage(buff)})
}

func (b *BoltServer[K, V]) HandleGetAllEcho(c echo.Context) error {
	if boltStorer, ok := b.Storage.(BoltStorer[int, *User]); ok {
		users, err := boltStorer.GetAll()
		if err != nil {
			log.Println("Error fetching all users:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "Error fetching users"})
		}
		return c.JSON(http.StatusOK, users)
	}

	log.Println("Error casting Storage to BoltStorer")
	return c.JSON(http.StatusInternalServerError, map[string]string{"msg": "Internal storage error"})
}
