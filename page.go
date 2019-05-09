package main

import "io/ioutil"

// Page - type representing a single editable page (from golang wiki.go tutorial)
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile("pages/"+filename, p.Body, 0600)
}
