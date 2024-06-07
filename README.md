
```
   ___  ____  ____  _______  __
  / _ \/ __ \/ __ \/_  __/ |/_/
 / , _/ /_/ / /_/ / / / _>  <  
/_/|_|\____/\____/ /_/ /_/|_|  
                               
```

# Introduction
Rootx is a command line tool designed for simplifying administrative tasks by allowing users to execute commands with root privileges without the need to switch users. It offers a range of features including module creation, database migration, seeder generation, and more, making it ideal for automation and management of projects.
                               




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


