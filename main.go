package rootx

import (
	"fmt"
	"os"
	"strings"

	"github.com/JubaerHossain/rootx/pkg/create"
	"github.com/spf13/cobra"
)

var asciiArt = `
___  ____  ____  _______  __
/ _ \/ __ \/ __ \/_  __/ |/_/
/ , _/ /_/ / /_/ / / / _>  <  
/_/|_|\____/\____/ /_/ /_/|_|  
 
A CLI tool for building RESTful APIs with Go
`

// Create a Cobra command for rootx
var rootCmd = &cobra.Command{
	Use:   "rootx",
	Short: "Rootx CLI Tool",
	Long:  asciiArt,
	Run: func(cmd *cobra.Command, args []string) {
		showMenu()
	},
}

// Convert hex color to RGB
func hexToRGB(hex string) (int, int, int) {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

// colorize function to add RGB colors
func colorize(text, hex string) string {
	r, g, b := hexToRGB(hex)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}

// showMenu displays the menu and handles user input
func showMenu() {
	for {
		fmt.Println(colorize(asciiArt, "#00FFFF"))                          // Cyan color for ASCII art
		fmt.Println(colorize("Select an option:", "#800080"))               // Magenta color for prompt
		fmt.Println(colorize("1. Create Module", "#FFFF00"))                 // Yellow color for option 7
		fmt.Println(colorize("2. Create Module with run", "#0080ff"))                // Green color for option 1
		fmt.Println(colorize("3. Create Migration", "#FFFF00"))             // Yellow color for option 2
		fmt.Println(colorize("4. Create Seeder", "#0000FF"))                // Blue color for option 3
		fmt.Println(colorize("5. Create Migration with Seeder", "#ccff66")) // Magenta color for option 4
		fmt.Println(colorize("6. Apply Migrations", "#00FFFF"))             // Cyan color for option 5
		fmt.Println(colorize("7. Run Seeders", "#0080ff"))                  // Green color for option 6
		fmt.Println(colorize("8. Run API Docs", "#FFFF00"))                 // Yellow color for option 7
		fmt.Println(colorize("0. Exit", "#FF0000"))          // Red color for return option
		fmt.Print(colorize("Enter the command number: ", "#006600"))        // Green color for the input prompt

		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println(colorize("Invalid input - please enter a number", "#FF0000")) // Red color for invalid input
			continue
		}

		if choice == 0 {
			fmt.Println(colorize("Bye!	👋🏼", "#FF0000"))
			os.Exit(0)
		}

		handleChoice(choice)
	}
}

func Yes() {
	for {
		fmt.Println(colorize("Did you migrate & seed?", "#FFA500")) // Orange color for a cautionary question
		fmt.Println(colorize("If not, please migrate & seed first.", "#FFA500")) // Orange color for guidance
		fmt.Println(colorize("Did you configure the Swagger docs in the /cmd/server/main.go file?", "#FFA500")) // Orange color for guidance
		fmt.Println(colorize("If not, please configure the Swagger docs in the /cmd/server/main.go file.", "#FFA500")) // Orange color for guidance
		fmt.Println(colorize("Follow the instructions in the README.md file.", "#FFA500")) // Orange color for guidance
		fmt.Println(colorize("1. Yes", "#00FF00")) // Green color for positive confirmation
		fmt.Println(colorize("2. No", "#FF0000")) // Red color for negative response
		fmt.Print(colorize("Do you want to continue? enter 1 or 2: ", "#006600")) // Orange color for the input prompt
		
		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println(colorize("Invalid input - please enter a number", "#FF0000")) // Red color for invalid input
			continue
		}
		if choice == 2 {
			showMenu()
		}
		if choice == 1 {
			if err := create.RunApiDocs(nil, nil); err != nil {
				fmt.Println()
				fmt.Println(colorize("Error running API docs", "#FF0000")) // Red color for error message
				fmt.Println()
				showMenu()
			}
		}
		break
	}
}


// handleChoice handles the user's choice
func handleChoice(choice int) {
	switch choice {
	case 1:
		moduleName := createMenu("Enter module name: ", "Module")
		args := []string{"create", moduleName}
		if err := create.Module(nil, args); err != nil {
			fmt.Println()
			fmt.Println("\x1b[31mError creating module\x1b[0m")
			fmt.Println()
			showMenu()
		}
	case 2:
		moduleName := createMenu("Enter module name: ", "Module")
		args := []string{"create", moduleName}
		if err := create.Run(nil, args); err != nil {
			fmt.Println()
			fmt.Println("\x1b[31mError creating module\x1b[0m")
			fmt.Println()
			showMenu()
		}
	case 3:
		migrationName := createMenu("Enter migration name: ", "Migration")
		args := []string{"create", migrationName}
		if err := create.MigrationCreate(nil, args); err != nil {
			fmt.Println()
			fmt.Println("\x1b[31mError creating migration\x1b[0m")
			fmt.Println()
			showMenu()
		}
	case 4:
		seederName := createMenu("Enter seeder name: ", "Seeder")
		args := []string{"create", seederName}
		if err := create.SeederCreate(nil, args); err != nil {
			fmt.Println()
			fmt.Println("\x1b[31mError creating seeder\x1b[0m")
			fmt.Println()
			showMenu()
		}
	case 5:
		name := createMenu("Enter name for migration and seeder: ", "Seeder")
		args := []string{"create", name}
		if err := create.MigrationWithSeederCreate(nil, args); err != nil {
			fmt.Println()
			fmt.Println("\x1b[31mError creating migration and seeder\x1b[0m")
			fmt.Println()
			showMenu()
		}
	case 6:
		if err := create.ApplyMigrations(nil, nil); err != nil {
			fmt.Println()
			fmt.Println(colorize(err.Error(), "#FF0000"))
			fmt.Println()
			showMenu()
		}
	case 7:
		if err := create.RunSeeders(nil, nil); err != nil {
			fmt.Println()
			fmt.Println(colorize(err.Error(), "#FF0000"))
			fmt.Println()
			showMenu()
		}
	case 8:
		Yes()
		
	default:
		fmt.Println(colorize("Invalid command", "#FF0000")) // Red color for invalid command
	}
}

func createMenu(name string, content string) string {
	moduleName := getUserInput(name)
	if len(moduleName) < 3 {
		fmt.Println()
		fmt.Printf("\x1b[31mError: %s name must be at least 3 characters long\x1b[0m\n", content)
		fmt.Println()
		return createMenu(name, content) // Re-prompt for module name
	}

	if strings.Contains(moduleName, " ") {
		fmt.Println()
		fmt.Printf("\x1b[31mError: %s name should not contain spaces\x1b[0m\n", content)
		fmt.Println()
		return createMenu(name, content) // Re-prompt for module name
	}

	return moduleName
}

// getUserInput prompts the user for input
func getUserInput(prompt string) string {
	fmt.Print(colorize(prompt, "#008000")) // Green color for user prompt
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

// Run is the main function to execute the root command
func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf(colorize("Execute error: %s", "#FF0000")+"\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(create.Create)
}
