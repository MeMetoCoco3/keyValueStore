package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var IssueList = template.
	Must(template.New("issueList").
		Parse(`
<h1>Users Table</h1>
	<table>
		<tr style='text-align: left'>
			<th>#</th>
			<th>Name</th>
			<th>Age</th>
			<th>Role</th>
		</tr>
	{{range .}}
		<tr>
			<td>{{.ID}}</td>
			<td><b>{{.Name}}</b></td>
			<td>{{.Age}}</td>
			<td>{{.Role}}</td>
		</tr>
	{{end}}
	</table>
		`))

type Server struct {
	Storage    Storer[int, *User]
	ListenAddr string
}

func NewServer(lAddr string) *Server {
	return &Server{
		Storage:    NewKVStore[int, *User](),
		ListenAddr: lAddr,
	}
}

func (s *Server) StartServer() {
	fmt.Printf("We are now listening on port %s", s.ListenAddr)
	http.HandleFunc("/put", s.HandlePut)
	log.Fatal(http.ListenAndServe(s.ListenAddr, nil))
}

func (s *Server) HandlePut(w http.ResponseWriter, r *http.Request) {
	var buff bytes.Buffer

	list, err := s.Storage.Iter(1, 2, 3, 6, 7)
	if err != nil {
		log.Fatalf("We got error %v", err)
	}

	err = IssueList.Execute(&buff, list)
	if err != nil {
		log.Fatalf("We got error %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buff.Bytes())
}
