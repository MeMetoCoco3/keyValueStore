package main

import "fmt"

func main() {
	s := NewKVStore[int, string]()

	s.Put(43, "Fortythree")

	v, err := s.Get(43)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)

	val, err := s.Delete(43)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(val)

}
