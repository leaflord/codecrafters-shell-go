package main

func main() {
	// fd := int(os.Stdin.Fd())
	// if !term.IsTerminal(fd) {
	// 	fmt.Println("Error: This must be run in a fully interactive terminal.")
	// 	return
	// }
	// Turn on raw mode to block terminal echo
	// oldState, err := term.MakeRaw(fd)
	// if err != nil {
	// 	fmt.Printf("Error entering raw mode: %v\n", err)
	// 	return
	// }
	// defer term.Restore(fd, oldState)
	console := NewConsole()
	console.Start()
}
