package cli

import (
	"GRIMIDIE/handler"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
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
		fmt.Println("5. Add Book")
		fmt.Println("6. Delete Books")
		fmt.Println("7. Create Admin")
		fmt.Println("8. Exit")
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
			cli.AddBook()
		case 6:
			cli.listBooks()
			cli.DeleteBook()
		case 7:
			cli.listUsersNotAdmin()
			//cli.showAddMenu()
		case 8:
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

	for {
		fmt.Print("Enter Name: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if isAlpha(name) {
			break
		} else {
			fmt.Println("Invalid input. Only letters are allowed.")
		}
	}
	for {
		fmt.Print("Enter Email:")
		fmt.Scanln(&email)
		if isValidEmail(email) {
			break
		} else {
			fmt.Println("Invalid email format. Please try again.")
		}
	}
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
	for {
		fmt.Print("Email:")
		fmt.Scanln(&email)
		if isValidEmail(email) {
			break
		} else {
			fmt.Println("Invalid email format. Please try again.")
		}
	}
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
		log.Print("Error listing Books: ", err)
		log.Fatal(err)
	}
	fmt.Println("Books listed successfully")
}

// Handler interface Addbook
func (cli *CLI) AddBook() {
	var judul, pengarang, tanggalTerbit string
	var qty int
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Book Title: ")
	judul, _ = reader.ReadString('\n')
	fmt.Println("Writer: ")
	reader = bufio.NewReader(os.Stdin)
	pengarang, _ = reader.ReadString('\n')
	fmt.Println("Publich Date: ")
	fmt.Scan(&tanggalTerbit)
	fmt.Println("Quantity: ")
	fmt.Scan(&qty)
	err := cli.Handler.AddBook(strings.ReplaceAll(judul, "\n", ""), strings.ReplaceAll(pengarang, "\n", ""), tanggalTerbit, qty)
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}
	fmt.Println("Add book successfully...")
}

func (cli *CLI) DeleteBook() {

	var BookID int
	fmt.Print("Which book: ")
	fmt.Scanln(&BookID)

	err := cli.Handler.DeleteBook(BookID)
	if err != nil {
		log.Print("Error Deleteing book: ", err)
		log.Fatal(err)
	}
}

// Handler interface CreatePinjam
func (cli *CLI) lendBook() {
	var choice int
	fmt.Println("Choose Book: ")
	fmt.Scan(&choice)
	err := cli.Handler.CreatePinjam(cli.CurrentUserID, choice, 1)
	if err != nil {
		log.Print("Error listing lending book: ", err)
		log.Fatal(err)
	}
	fmt.Println("Displaying Books List...")
}

// Handler interface listPeminjaman
func (cli *CLI) listPinjam() {
	err := cli.Handler.ListPeminjaman(cli.CurrentUserID)
	if err != nil {
		log.Print("Error listing pinjam: ", err)
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

func (cli *CLI) listUsersNotAdmin() {
	err := cli.Handler.ListUsersNotAdmin()
	if err != nil {
		log.Print("Error listing users: ", err)
		log.Fatal(err)
	}

	var choice int
	fmt.Print("Choose the ID of the user to grant admin rights: ")
	fmt.Scan(&choice)

	err = cli.Handler.UpdateUserAdminStatus(choice, true)
	if err != nil {
		log.Print("Error updating user status: ", err)
		log.Fatal(err)
	}

	fmt.Println("User granted admin rights successfully.")
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

// Region Func Validasi Input

func isAlpha(input string) bool {
	for _, letter := range input {
		if !unicode.IsLetter(letter) && letter != ' ' {
			return false
		}
	}
	return true
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
