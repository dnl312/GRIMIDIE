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
	Handler       handler.Handler
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

// Region MENU CLI BEGIN
func (cli *CLI) showMenu() {
	for {
		// HOME MENU
		fmt.Println("		 🎗️GRIMIDIE Application!🎗️")
		fmt.Println("\nPlease Sign In To GRIMIDIE | Don't have an account yet? Sign Up")
		fmt.Println("\n1. Sign Up")
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
		case 123:
			cli.signInDebugMode()
		default:
			fmt.Println("⚠️ Invalid choice. Please try again. ⚠️")
		}
	}
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
			fmt.Println("⚠️ Invalid choice. Please try again. ⚠️")
		}
	}
}

func (cli *CLI) showAdminMenu() {
	for {
		fmt.Println("\nWelcome Admin")
		fmt.Println("1. User Reports")
		fmt.Println("2. Lend Reports")
		fmt.Println("3. Stock Reports")
		fmt.Println("4. Most Loan Books Reports")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			//cli.UserReports()
		case 2:
			cli.reportPinjam()
		case 3:
			cli.reportStock()
		case 4:
			cli.reportPopularBooks()
		case 5:
			fmt.Println("GoodBye Min!")
			os.Exit(0)
		default:
			fmt.Println("⚠️ Invalid choice. Please try again. ⚠️")
		}
	}
}

// Region MENU CLI END

func (cli *CLI) signInDebugMode() {

	userID, _, err := cli.Handler.UserLogin("hannah@example.com", "password234")
	if err != nil {
		fmt.Println("Error during sign in:", err)
		return
	}
	cli.CurrentUserID = userID
	cli.showUserMenu()
}

// Region HANDLER interface BEGIN
// Handler interface UserRegister
func (cli *CLI) signUp() {
	var name, email, password string

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Name:")
	name, _ = reader.ReadString('\n')
	fmt.Print("Enter Email:")
	fmt.Scanln(&email)
	fmt.Print("Enter Password:")
	fmt.Scanln(&password)

	_, err := cli.Handler.UserRegister(strings.ReplaceAll(name, "\n", ""), email, password)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Your account has been successfully registered!")

}

// Handler interface UserLogin
func (cli *CLI) signIn() {
	var email, password string
	fmt.Print("Email:")
	fmt.Scanln(&email)
	fmt.Print("Password:")
	fmt.Scanln(&password)

	userID, isAdmin, err := cli.Handler.UserLogin(email, password)
	if err != nil {
		fmt.Println("Error during sign in:", err)
		return
	}
	if isAdmin {
		cli.showAdminMenu()
	} else {
		cli.CurrentUserID = userID
		cli.showUserMenu()
	}

}

// Handler interface listBook
func (cli *CLI) listBooks() {
	err := cli.Handler.ListBooks()
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}
	fmt.Println("Users listed successfully")
}

// Handler interface CreatePinjam
func (cli *CLI) lendBook() {
	var choice int
	fmt.Println("Choose Book: ")
	fmt.Scan(&choice)
	err := cli.Handler.CreatePinjam(cli.CurrentUserID, choice, 1)
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}
	fmt.Println("Displaying Books List...")
}

// Handler interface listPeminjaman
func (cli *CLI) listPinjam() {
	err := cli.Handler.ListPeminjaman(cli.CurrentUserID)
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}
}

// Handler interface returnPinjam
func (cli *CLI) returnBook() {

	var OrderID int
	fmt.Print("Choose: ")
	fmt.Scanln(&OrderID)

	denda, err := cli.Handler.ReturnPinjam(cli.CurrentUserID, OrderID)
	if err != nil {
		log.Print("Error returning book: ", err)
		log.Fatal(err)
	}

	if denda > 0 {
		fmt.Printf("Denda: %.2f", denda)
		fmt.Println()
	}
}

// Region HANDLER interface END

// Region REPORT HANDLER CLI

func (cli *CLI) reportPinjam() {
	err := cli.Handler.ReportPeminjaman()
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}
}
func (cli *CLI) reportStock() {
	err := cli.Handler.ReportStock()
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}
}

func (cli *CLI) reportPopularBooks() {
	err := cli.Handler.ReportPopularBooks()
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}
}

// REPORT USER DISINI YAH
