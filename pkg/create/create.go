package create

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"github.com/schollz/progressbar/v3"
	"github.com/gertd/go-pluralize"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	AppRoot = "domain"

	TemplateDir = "../../template"

	ServiceDir     = "service"
	EntityDir      = "entity"
	RepositoryDir  = "repository"
	PersistenceDir = "infrastructure/persistence"
	http           = "infrastructure/transport/http"
)

var AppName string

var Create = &cobra.Command{
	Use:  "create",
	Args: cobra.MinimumNArgs(2),
	RunE: Run,
}

func CreateProgressBar(description string) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(100,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetDescription(fmt.Sprintf("[magenta]%s: ", description)),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(100),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[magenta]=[reset]",
			SaucerPadding: " ",
			SaucerHead:    "[magenta]>[reset]",
		}),
	)

	return bar
}


func Run(cmd *cobra.Command, args []string) error {

	bar := CreateProgressBar("Creating module: ")

	// Perform your tasks here
	for i := 0; i <= 100; i++ {
		// Update progress bar
		bar.Add(1)
		time.Sleep(100 * time.Millisecond) // Simulate some work being done
	}

	if len(args) < 2 {
		return errors.New("not enough arguments")
	}
	AppName = "github.com/JubaerHossain/rootx"
	name := args[1]
	name = Lower(Plural(name))
	fs := afero.NewBasePathFs(afero.NewOsFs(), AppRoot+"/")
	if err := createFolders(fs, name); err != nil {
		return err
	}
	if err := createFiles(fs, name); err != nil {
		return err
	}
	if err := MigrationWithSeederCreate(nil, args); err != nil {
		return err
	}
	if err := RunApp(nil, nil); err != nil {
		return err
	}
	return nil
}

func getModuleName() (string, error) {
	modFile, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer modFile.Close()

	scanner := bufio.NewScanner(modFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("module name not found in go.mod")
}

func createFolders(fs afero.Fs, name string) error {
	fs.Mkdir(name, 0755)
	dirs := []string{ServiceDir, EntityDir, RepositoryDir, PersistenceDir, http}
	for _, dir := range dirs {
		if err := fs.MkdirAll(path.Join(name, dir), 0755); err != nil {
			return err
		}
	}
	return nil
}
func createSingleFolders(fs afero.Fs, name string) error {
	fs.Mkdir(name, 0755)
	return nil
}

func createFiles(fs afero.Fs, name string) error {
	createFile(fs, name, path.Join(TemplateDir, "service.stub"), path.Join(name, ServiceDir, name+".go"))
	createFile(fs, name, path.Join(TemplateDir, "entity.stub"), path.Join(name, EntityDir, name+".go"))
	createFile(fs, name, path.Join(TemplateDir, "repository.stub"), path.Join(name, RepositoryDir, name+".go"))
	createFile(fs, name, path.Join(TemplateDir, "persistence.stub"), path.Join(name, PersistenceDir, name+".go"))
	createFile(fs, name, path.Join(TemplateDir, "handler.stub"), path.Join(name, http, "handler.go"))
	createFile(fs, name, path.Join(TemplateDir, "route.stub"), path.Join(name, http, "route.go"))

	return nil
}

func createFile(fs afero.Fs, name, stubPath, filePath string) error {
	fs.Create(filePath)

	_, filename, _, _ := runtime.Caller(1)
	stubPath = path.Join(path.Dir(filename), stubPath)

	contents, err := fileContents(stubPath)
	if err != nil {
		return err
	}
	contents = replaceStub(contents, name)

	if err := overwrite(AppRoot+"/"+filePath, contents); err != nil {
		return err
	}
	return nil
}

func fileContents(file string) (string, error) {
	a := afero.NewOsFs()
	contents, err := afero.ReadFile(a, file)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func overwrite(file string, message string) error {
	a := afero.NewOsFs()
	return afero.WriteFile(a, file, []byte(message), 0666)
}

func replaceStub(content string, name string) string {

	content = strings.Replace(content, "{{TitleName}}", Title(name), -1)
	content = strings.Replace(content, "{{PluralLowerName}}", Lower(Plural(name)), -1)
	content = strings.Replace(content, "{{SingularLowerName}}", Lower(Singular(name)), -1)
	content = strings.Replace(content, "{{SingularCapitalName}}", UpperCamelCase(Singular(name)), -1)
	content = strings.Replace(content, "{{PluralCapitalName}}", UpperCamelCase(Plural(name)), -1)
	content = strings.Replace(content, "{{AppName}}", AppName, -1)
	content = strings.Replace(content, "{{AppRoot}}", AppRoot, -1)
	return content
}

func Plural(name string) string {
	pluralize := pluralize.NewClient()
	return pluralize.Plural(name)
}

func Singular(name string) string {
	pluralize := pluralize.NewClient()
	return pluralize.Singular(name)
}

func Lower(name string) string {
	return strings.ToLower(name)
}

func Title(name string) string {
	return strings.Title(Lower(name))
}
func UpperCamelCase(name string) string {
	return strings.Title(name)
}

func MigrationCreate(cmd *cobra.Command, args []string) error {
	bar := CreateProgressBar("Creating migration: ")
	// Perform your tasks here
	for i := 0; i <= 100; i++ {
		// Update progress bar
		bar.Add(1)
		time.Sleep(100 * time.Millisecond) // Simulate some work being done
	}

	if len(args) < 2 {
		return errors.New("not enough arguments")
	}
	name := args[1]
	name = Lower(Plural(name))
	if err := createMigrationFile(name); err != nil {
		fmt.Print(err)
		return errors.New("error creating migration file")
	}

	return nil
}

func createMigrationFile(name string) error {
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		os.Mkdir("migrations", 0755)
	}
	timestamp := time.Now().Format("2006_01_02_150405")
	filename := filepath.Join("migrations", fmt.Sprintf("%s_%s.sql", timestamp, name))
	content := fmt.Sprintf("-- Migration %s\n\n", name) +
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", name) +
		"    id SERIAL PRIMARY KEY,\n" +
		"    name VARCHAR(100) NOT NULL,\n" +
		"    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
		"    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP\n" +
		");\n\n" +
		fmt.Sprintf("CREATE INDEX ON %s (name);\n", name) // Modify column_name with the actual column name

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}
	return nil
}

func SeederCreate(cmd *cobra.Command, args []string) error {
	bar := CreateProgressBar("Creating seeder: ")
	// Perform your tasks here
	for i := 0; i <= 100; i++ {
		// Update progress bar
		bar.Add(1)
		time.Sleep(100 * time.Millisecond) // Simulate some work being done
	}

	if len(args) < 2 {
		return errors.New("not enough arguments")
	}
	name := args[1]
	name = Lower(Plural(name))
	if err := createSeedFile(name); err != nil {
		return errors.New("error creating seeder file")
	}

	return nil
}

func createSeedFile(tableName string) error {
	if _, err := os.Stat("seeds"); os.IsNotExist(err) {
		os.Mkdir("seeds", 0755)
	}
	timestamp := time.Now().Format("2006_01_02_150405")
	filename := filepath.Join("seeds", fmt.Sprintf("%s_%s_seeder.sql", timestamp, tableName))
	content := fmt.Sprintf("-- Seeder for table %s\n\n", tableName) +
		fmt.Sprintf("INSERT INTO %s (name, created_at, updated_at) VALUES\n", tableName) +
		"    ('Value1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),\n" +
		"    ('Value2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);\n"

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create seed file: %w", err)
	}
	return nil
}

func MigrationWithSeederCreate(cmd *cobra.Command, args []string) error {
	bar := CreateProgressBar("Creating migration and seeder: ")
	// Perform your tasks here
	for i := 0; i <= 100; i++ {
		// Update progress bar
		bar.Add(1)
		time.Sleep(100 * time.Millisecond) // Simulate some work being done
	}

	if len(args) < 2 {
		return errors.New("not enough arguments")
	}
	name := args[1]
	name = Lower(Plural(name))
	if err := createMigrationFile(name); err != nil {
		return errors.New("error creating migration file")
	}

	if err := createSeedFile(name); err != nil {
		return errors.New("error creating seeder file")
	}

	return nil
}

func connectDB(dbUser string, dbPassword string, dbHost string, dbPortStr string, dbName string) (*pgxpool.Pool, error) {

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert DB_PORT to int: %w", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	config.MaxConnIdleTime = 10 * time.Minute
	config.MaxConnLifetime = 60 * time.Minute // Set to 1 hour
	config.MaxConns = 50000                    // Adjust based on your environment
	config.MinConns = 100

	return pgxpool.NewWithConfig(context.Background(), config)
}

func getUserInput(prompt string) string {
	// ANSI escape code for green color
	green := "\033[32m"
	// ANSI escape code to reset color
	reset := "\033[0m"

	// Print prompt in green color
	fmt.Print(green + prompt + reset)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func ApplyMigrations(cmd *cobra.Command, args []string) error {
	
	dbUser := getUserInput("Enter database user: ")
	dbPassword := getUserInput("Enter database password: ")
	dbHost := getUserInput("Enter database host: ")
	dbPortStr := getUserInput("Enter database port: ")
	dbName := getUserInput("Enter database name: ")

	pool, err := connectDB(dbUser, dbPassword, dbHost, dbPortStr, dbName)
	if err != nil {
		return fmt.Errorf("failed to connect to database")
	}
	defer pool.Close()

	migrationsDir := "migrations"
	if err := executeScriptsInDirectory(pool, migrationsDir); err != nil {
		return fmt.Errorf("failed to execute migration scripts")
	}

	bar := CreateProgressBar("migration: ")
	// Perform your tasks here
	for i := 0; i <= 100; i++ {
		// Update progress bar
		bar.Add(1)
		time.Sleep(100 * time.Millisecond) // Simulate some work being done
	}



	return nil
}

func RunSeeders(cmd *cobra.Command, args []string) error {
	bar := CreateProgressBar("Seeding: ")
	// Perform your tasks here
	for i := 0; i <= 100; i++ {
		// Update progress bar
		bar.Add(1)
		time.Sleep(100 * time.Millisecond) // Simulate some work being done
	}

	dbUser := getUserInput("Enter database user: ")
	dbPassword := getUserInput("Enter database password: ")
	dbHost := getUserInput("Enter database host: ")
	dbPortStr := getUserInput("Enter database port: ")
	dbName := getUserInput("Enter database name: ")

	pool, err := connectDB(dbUser, dbPassword, dbHost, dbPortStr, dbName)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	seedsDir := "seeds"
	if err := executeScriptsInDirectory(pool, seedsDir); err != nil {
		return fmt.Errorf("failed to execute seeder scripts: %w", err)
	}
	return nil
}

func executeScriptsInDirectory(pool *pgxpool.Pool, directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(directory, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		_, err = pool.Exec(context.Background(), string(content))
		if err != nil {
			return fmt.Errorf("failed to execute file %s: %w", filePath, err)
		}
	}
	return nil
}
func createDocsFile(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		if err := os.Mkdir(name, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	filename := filepath.Join(name, "docs.go")
	content := fmt.Sprintf("package %s\n\n", name) +
		"// @title RootX API\n" +
		"// @description This is the API documentation for RootX\n" +
		"// @version 1\n" +
		"// @host localhost:8080\n" +
		"// @BasePath /\n"

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create docs file: %w", err)
	}
	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}
	return nil
}

func createMainFile(templatePath, targetPath string) error {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to read template file: %w", err)
		}
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to create main.go file: %w", err)
		}
	}
	return nil
}

func createEnvFile() error {
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		templateFile := filepath.Join("template", "env.stub")
		content, err := os.ReadFile(templateFile)
		if err != nil {
			return fmt.Errorf("failed to read template file: %w", err)
		}
		if err := os.WriteFile(envFile, content, 0644); err != nil {
			return fmt.Errorf("failed to create .env file: %w", err)
		}
	}
	return nil
}

func RunApp(cmd *cobra.Command, args []string) error {
	bar := CreateProgressBar("App Running: ")
	// Perform your tasks here
	for i := 0; i <= 100; i++ {
		// Update progress bar
		bar.Add(1)
		time.Sleep(100 * time.Millisecond) // Simulate some work being done
	}

	if err := createEnvFile(); err != nil {
		return fmt.Errorf("error creating .env file: %w", err)
	}

	templatePath := "template/main.stub"
	targetPath := "./cmd/server/main.go"
	if err := createMainFile(templatePath, targetPath); err != nil {
		return fmt.Errorf("error creating main.go file: %w", err)
	}

	if err := runCommand("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	// Run the server
	if err := runServer(); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
func RunApiDocs(cmd *cobra.Command, args []string) error {
	bar := CreateProgressBar("API Docs Generating: ")
	// Perform your tasks here
	for i := 0; i <= 100; i++ {
		// Update progress bar
		bar.Add(1)
		time.Sleep(100 * time.Millisecond) // Simulate some work being done
	}
	if err := createDocsFile("docs"); err != nil {
		return fmt.Errorf("error creating docs file: %w", err)
	}
	
	if err := runCommand("go", "get", "github.com/swaggo/swag/cmd/swag@v1.16.3"); err != nil {
		return fmt.Errorf("failed to install swag: %w", err)
	}

	if err := runCommand("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	if err := runCommand("swag", "init", "-g", "./cmd/server/main.go"); err != nil {
		return fmt.Errorf("failed to run swag init: %w", err)
	}

	// Run the server
	if err := runServer(); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}

// Helper function to terminate a running process
func terminateRunningProcess() error {
	pidFile := "server.pid"
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return nil // No running process found
	}

	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		return fmt.Errorf("failed to read pid file: %w", err)
	}

	var pid int
	_, err = fmt.Sscanf(string(pidData), "%d", &pid)
	if err != nil {
		return fmt.Errorf("failed to parse pid: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to terminate process: %w", err)
	}

	// Wait for the process to terminate
	time.Sleep(2 * time.Second)

	if err := os.Remove(pidFile); err != nil {
		return fmt.Errorf("failed to remove pid file: %w", err)
	}

	return nil
}

// Helper function to run the server and save its PID
func runServer() error {
	cmd := exec.Command("go", "run", "./cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	pidFile := "server.pid"
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644); err != nil {
		return fmt.Errorf("failed to write pid file: %w", err)
	}

	return cmd.Wait()
}
