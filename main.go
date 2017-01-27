package main

var h *Howdoi

func init() {
	h = Init()
}

func main() {
	h.Execute()
}
