package main

import (
	"log"
)

func main() {
	s, err := NewBoltServer("Bunny.db", ":1337")
	if err != nil {
		log.Fatalln(err)
	}
	u1 := NewUser("Vidal", "Macho", 28)
	u2 := NewUser("Alice", "Fierce", 25)
	u3 := NewUser("Bob", "Brave", 30)
	u4 := NewUser("Charlie", "Swift", 22)
	u5 := NewUser("Diana", "Strong", 27)
	u6 := NewUser("Ethan", "Clever", 35)
	u7 := NewUser("Fiona", "Bold", 26)
	u8 := NewUser("George", "Quick", 29)
	u9 := NewUser("Hannah", "Kind", 31)
	u10 := NewUser("Isaac", "Loyal", 24)
	u11 := NewUser("Julia", "Smart", 28)
	u12 := NewUser("Kevin", "Sharp", 33)
	u13 := NewUser("Laura", "Fearless", 21)
	u14 := NewUser("Mike", "Steady", 34)
	u15 := NewUser("Nina", "Happy", 23)
	u16 := NewUser("Oscar", "Calm", 32)
	u17 := NewUser("Paula", "Wise", 26)
	u18 := NewUser("Quentin", "Energetic", 29)
	u19 := NewUser("Rachel", "Brilliant", 30)
	u20 := NewUser("Sam", "Resilient", 27)
	err = s.Storage.PutB(1, u1)
	log.Println(err)
	err = s.Storage.PutB(2, u2)

	log.Println(err)
	err = s.Storage.PutB(3, u3)

	log.Println(err)
	s.Storage.PutB(4, u4)
	s.Storage.PutB(5, u5)
	s.Storage.PutB(6, u6)
	s.Storage.PutB(7, u7)
	s.Storage.PutB(8, u8)
	s.Storage.PutB(9, u9)
	s.Storage.PutB(10, u10)
	s.Storage.PutB(11, u11)
	s.Storage.PutB(12, u12)
	s.Storage.PutB(13, u13)
	s.Storage.PutB(14, u14)
	s.Storage.PutB(15, u15)
	s.Storage.PutB(16, u16)
	s.Storage.PutB(17, u17)
	s.Storage.PutB(18, u18)
	s.Storage.PutB(19, u19)
	s.Storage.PutB(20, u20)
	v, err := s.Storage.GetB(17)
	log.Println(v)
	log.Println(err)
	data, err := s.Storage.GetAll()

	for k, v := range data {
		log.Println("GETALL:", k, v)

	}
}
