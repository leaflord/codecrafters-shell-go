package main

func main() {
	console := NewConsole()
	defer console.Clean()
	console.Start()
}
