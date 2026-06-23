package main

func main() {
	console := NewConsole()
	defer console.Clean()
	console.Start(func(r rune) bool {
		return handleRune(console, r)
	})
}
