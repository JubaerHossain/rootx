
```
   ___  ____  ____  _______  __
  / _ \/ __ \/ __ \/_  __/ |/_/
 / , _/ /_/ / /_/ / / / _>  <  
/_/|_|\____/\____/ /_/ /_/|_|  
                               
```

# Introduction
Rootx is a command line tool designed for simplifying administrative tasks by allowing users to execute commands with root privileges without the need to switch users. It offers a range of features including module creation, database migration, seeder generation, and more, making it ideal for automation and management of projects.
                               



## Installation

### macOS and Linux

1. Download the appropriate tar.gz file for your platform from the [releases page](https://github.com/JubaerHossain/rootx/releases).
2. Extract the tar.gz file.
3. Run the installation script.

```bash
# Download and extract the tar.gz file
curl -L -o rootx-darwin-amd64.tar.gz <URL to the tar.gz file>
tar -xvzf rootx-darwin-amd64.tar.gz
cd rootx-darwin-amd64

# Install
chmod +x install.sh
sudo ./install.sh
```

### Windows

    1. Download the `rootx-windows-amd64.zip` file from the [releases page](https://github.com/JubaerHossain/rootx/releases).
    2. Extract the zip file.
    3. Run the installation script.

```powershell
# Download and extract the zip file
Invoke-WebRequest -Uri <URL to the zip file> -OutFile rootx-windows-amd64.zip
Expand-Archive -Path rootx-windows-amd64.zip -DestinationPath rootx-windows-amd64
cd rootx-windows-amd64

# Install
.\install.ps1
```

## Usage

Once installed, you can use the `rootx` command in your terminal.

```bash
rootx
```

This will display the rootx CLI menu where you can choose various options to create modules, migrations, seeders, etc.

### CLI Options

1. **Create Module**
2. **Create Migration**
3. **Create Seeder**
4. **Create Migration with Seeder**
5. **Apply Migrations**
6. **Run Seeders**
7. **Generate API Documentation**
0. **Exit**

### Example Commands

#### Create a Module

```bash
rootx create module <moduleName>
```



## Or Installation by go get 
```bash
  go get -u github.com/JubaerHossain/rootx
```

## Run your app
```bash
   cmd
   ├── rootx
   │   ├── main.go
```

```bash
  package main

import (
	"github.com/JubaerHossain/rootx"
)

func main() {
	rootx.Run()
}

```

## Usage/Examples
```bash
  go run ./cmd/rootx
```
### create module
```bash
  ___  ____  ____  _______  __
/ _ \/ __ \/ __ \/_  __/ |/_/
/ , _/ /_/ / /_/ / / / _>  <  
/_/|_|\____/\____/ /_/ /_/|_|  
 
A CLI tool for building RESTful APIs with Go

Select an option:
1. Create Module
2. Create Migration
3. Create Seeder
4. Create Migration with Seeder
5. Apply Migrations
6. Run Seeders
7. Run API Docs
0. Return to Main Menu
Enter the command number: 1
Enter module name: users
```

### create migration
```bash
___  ____  ____  _______  __
/ _ \/ __ \/ __ \/_  __/ |/_/
/ , _/ /_/ / /_/ / / / _>  <  
/_/|_|\____/\____/ /_/ /_/|_|  
 
A CLI tool for building RESTful APIs with Go

Select an option:
1. Create Module
2. Create Migration
3. Create Seeder
4. Create Migration with Seeder
5. Apply Migrations
6. Run Seeders
7. Run API Docs
0. Return to Main Menu
Enter the command number: 2
Enter migration name: users
```

### create seeder
```bash
___  ____  ____  _______  __
/ _ \/ __ \/ __ \/_  __/ |/_/
/ , _/ /_/ / /_/ / / / _>  <  
/_/|_|\____/\____/ /_/ /_/|_|  
 
A CLI tool for building RESTful APIs with Go

Select an option:
1. Create Module
2. Create Migration
3. Create Seeder
4. Create Migration with Seeder
5. Apply Migrations
6. Run Seeders
7. Run API Docs
0. Return to Main Menu
Enter the command number: 3
Enter migration name: users
```

### create migration with seeder
```bash
___  ____  ____  _______  __
/ _ \/ __ \/ __ \/_  __/ |/_/
/ , _/ /_/ / /_/ / / / _>  <  
/_/|_|\____/\____/ /_/ /_/|_|  
 
A CLI tool for building RESTful APIs with Go

Select an option:
1. Create Module
2. Create Migration
3. Create Seeder
4. Create Migration with Seeder
5. Apply Migrations
6. Run Seeders
7. Run API Docs
0. Return to Main Menu
Enter the command number: 4
Enter migration name: users
```

### apply migrations
```bash

___  ____  ____  _______  __
/ _ \/ __ \/ __ \/_  __/ |/_/
/ , _/ /_/ / /_/ / / / _>  <  
/_/|_|\____/\____/ /_/ /_/|_|  
 
A CLI tool for building RESTful APIs with Go

Select an option:
1. Create Module
2. Create Migration
3. Create Seeder
4. Create Migration with Seeder
5. Apply Migrations
6. Run Seeders
7. Run API Docs
0. Return to Main Menu
Enter the command number: 5
Applying migrations...
Enter database user: postgres
Enter database password: password
Enter database host: localhost
Enter database port: 5433
Enter database name: starter_api

```

### run seeders
```bash

___  ____  ____  _______  __
/ _ \/ __ \/ __ \/_  __/ |/_/
/ , _/ /_/ / /_/ / / / _>  <  
/_/|_|\____/\____/ /_/ /_/|_|  
 
A CLI tool for building RESTful APIs with Go

Select an option:
1. Create Module
2. Create Migration
3. Create Seeder
4. Create Migration with Seeder
5. Apply Migrations
6. Run Seeders
7. Run API Docs
0. Return to Main Menu
Enter the command number: 6
Applying migrations...
Enter database user: postgres
Enter database password: password
Enter database host: localhost
Enter database port: 5433
Enter database name: starter_api
  
  ```

### generate api documentation
  #### Add the following code in your main.go file
```bash
package main

import (
	_ "github.com/JubaerHossain/rootx/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)


// rest of the code


func setupRoutes(application *app.App) http.Handler {

	// Register Swagger routes
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
  
  
  
}
```

  #### run the following command
```bash
  go run ./cmd/rootx
```

  - Select an option: 7



## Authors
- [@JubaerHossain](https://www.github.com/JubaerHossain)

## License
[MIT](https://choosealicense.com/licenses/mit/)


