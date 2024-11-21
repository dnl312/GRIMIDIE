package cli

import (
	"GRIMIDIE/handler"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type CLI struct {
	Handler handler.Handler
	CurrentUserID int
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
		fmt.Println("           Don't have an account yet? Sign Up")
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
			cli.signInDebugMode()
			//cli.signIn()
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
	fmt.Print("Enter Name:")
	name, _ = reader.ReadString('\n')
	fmt.Print("Enter Email:")
	fmt.Scanln(&email)
	fmt.Print("Enter Password:")
	fmt.Scanln(&password)

	_,err := cli.Handler.UserRegister(strings.ReplaceAll(name, "\n", ""), email, password)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Your account has been successfully registered!")

}

// Fungsi untuk proses sign in
func (cli *CLI) signIn() {
	var email, password string
	fmt.Print("Email:")
	fmt.Scanln(&email)
	fmt.Print("Password:")
	fmt.Scanln(&password)

	userID, err := cli.Handler.UserLogin(email, password)
	if err != nil {
		fmt.Println("Error during sign in:", err)
		return
	}

	cli.CurrentUserID = userID 
	cli.showUserMenu()
}

func (cli *CLI) signInDebugMode() {

	userID, err := cli.Handler.UserLogin("jack@example.com", "password890")
	if err != nil {
		fmt.Println("Error during sign in:", err)
		return
	}

	cli.CurrentUserID = userID 
	cli.showUserMenu()
}

func (cli *CLI) showUserMenu() {
	for {
		fmt.Println("\nHere is a list of books you can choose from.")
		cli.listBooks()
		fmt.Println("1. Lend Book")
		fmt.Println("2. Return Book")
		fmt.Println("3. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			// lend book from library
			cli.lendBook()
		case 2:
			// return book from library
			cli.listPinjam()
			cli.returnBook()
		case 3:
			fmt.Println("GoodBye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
func (cli *CLI) listBooks() {
	err := cli.Handler.ListBooks()
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}
	fmt.Println("Users listed successfully")
}

func (cli *CLI) lendBook() {
	// logic
	fmt.Println("Displaying Books List...")
}


func (c *CLI) listPinjam() {
	err := c.Handler.ListPeminjaman(c.CurrentUserID)
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}
}

func (cli *CLI) returnBook() {

	var OrderID int
	fmt.Print("Choose: ")
	fmt.Scanln(&OrderID)

	denda,err := cli.Handler.ReturnPinjam(OrderID)
	if err != nil {
		log.Print("Error returning book: ", err)
		log.Fatal(err)
	}

	if(denda>0){
		fmt.Printf("Denda: %.2f", denda)
		fmt.Println()
	}
}
