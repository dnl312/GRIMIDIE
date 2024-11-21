package cli

import (
	"GRIMIDIE/handler"
	"bufio"
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
		fmt.Println("Welcome to the GRIMIDIE Application! | Sign In To GRIMIDIE")
		fmt.Println("Don't have an account yet? Sign Up")
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

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter Name:")
	name, _ = reader.ReadString('\n')
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

// Fungsi untuk proses sign in
func (cli *CLI) signIn() {
	var email, password string
	fmt.Println("Email:")
	fmt.Scanln(&email)
	fmt.Println("Password:")
	fmt.Scanln(&password)

	if err := cli.Handler.UserLogin(email, password); err != nil {
		fmt.Println("Error during sign in:", err)
		return
	}
}
