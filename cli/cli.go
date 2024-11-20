package cli

import (
	"GRIMIDIE/handler"
	"fmt"
	"os"
)

type CLI struct {
	Handler handler.Handler
}

func NewCLI(handler handler.Handler) *CLI {
	return &CLI{
		Handler: handler,
	}
}

func (cli *CLI) Init() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	cli.showMenu()
}

func (cli *CLI) showMenu() {
	for {
		// Menampilkan menu pilihan
		fmt.Println("Welcome to the CLI Application!")
		fmt.Println("1. Sign Up")
		fmt.Println("2. Sign In")
		fmt.Println("3. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			// Sign Up
			cli.signUp()
		case 2:
			// Sign In
			cli.signIn()
		case 3:
			fmt.Println("GoodBye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// Fungsi untuk proses sign up
func (cli *CLI) signUp() {
	var name, email, password string
	fmt.Println("Enter Name:")
	fmt.Scanln(&name)
	fmt.Println("Enter Email:")
	fmt.Scanln(&email)
	fmt.Println("Enter Password:")
	fmt.Scanln(&password)

	err := cli.Handler.UserRegister(name, email, password)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Your account has been successfully registered!")
}
