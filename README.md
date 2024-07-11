
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
2. Create Module with run
3. Create Migration
4. Create Seeder
5. Create Migration with Seeder
6. Apply Migrations
7. Run Seeders
8. Scaffold Auth
9. Run API Docs
0. Exit
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
2. Create Module with run
3. Create Migration
4. Create Seeder
5. Create Migration with Seeder
6. Apply Migrations
7. Run Seeders
8. Scaffold Auth
9. Run API Docs
0. Exit
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
2. Create Module with run
3. Create Migration
4. Create Seeder
5. Create Migration with Seeder
6. Apply Migrations
7. Run Seeders
8. Scaffold Auth
9. Run API Docs
0. Exit
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
2. Create Module with run
3. Create Migration
4. Create Seeder
5. Create Migration with Seeder
6. Apply Migrations
7. Run Seeders
8. Scaffold Auth
9. Run API Docs
0. Exit
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
2. Create Module with run
3. Create Migration
4. Create Seeder
5. Create Migration with Seeder
6. Apply Migrations
7. Run Seeders
8. Scaffold Auth
9. Run API Docs
0. Exit
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
2. Create Module with run
3. Create Migration
4. Create Seeder
5. Create Migration with Seeder
6. Apply Migrations
7. Run Seeders
8. Scaffold Auth
9. Run API Docs
0. Exit
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

### File upload documentation

   - Add the following code in your codebase
   - Here **image** is the name of the form field and **foods** is the folder name
   - The uploaded file will be saved in the **uploads/foods** folder
   - You can use **local storage** or can use **S3 AWS**
   - set the following environment variable in the .env file
```bash
STORAGE_DISK=local // s3
STORAGE_PATH=storage
AWS_ACCESS_KEY=
AWS_SECRET_KEY=
AWS_REGION=ap-southeast-1
AWS_BUCKET=aws-bucket
AWS_ENDPOINT=https://s3.ap-southeast-1.amazonaws.com
```
   - Example:

```bash
	imageMetadata, err := r.app.FileUpload.FileUpload(req, "image", "foods")  // "image" is the name of the form field and "foods" is the folder name
	if err != nil {
		return fmt.Errorf("failed to upload image: %v", err)
	}
```


### Create Auth Scaffold 
    - Example:
```bash
  go run ./cmd/rootx
```
```bash
Select an option:
1. Create Module
2. Create Module with run
3. Create Migration
4. Create Seeder
5. Create Migration with Seeder
6. Apply Migrations
7. Run Seeders
8. Scaffold Auth
9. Run API Docs
0. Exit
Enter the command number: 8
```
### Auth
    - Add the following code in your route.go
    - will protected the routes
    - Example :

```bash
router.Handle("GET /employees", middleware.LimiterMiddleware(middleware.AuthMiddleware(http.HandlerFunc(handler.GetEmployees))))
```
### Database Connection
    - Added postgres and mysql database connection
    - By default, it will connect to the postgres database
    - You can change the database connection by changing the database connection in .env file
    - Example:
```bash
DB_TYPE=postgres // or mysql
DB_HOST=db_host
DB_PORT=db_port
DB_NAME=db_name
DB_USER=db_user
DB_PASSWORD=db_password
```




## Authors
- [@JubaerHossain](https://www.github.com/JubaerHossain)

## License
[MIT](LICENSE)


