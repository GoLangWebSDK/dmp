# Data Management Package

This module provides adapters for various types of database connections, and exposes interfaces for performing various operations on the stored data, like CRUD operations, migrations, pagination, etc.

It has one sub-package called "database" which has the top-level interface exposing all the operations you need to do to set up a database connection, migrate and seed the data, and expose the connection to the repository layer.

```
// The new database configuration expects a configuration struct...
config := &database.DBConfig{
    DBName: "sdk-dev-db",
    DBUser: "sdk-admin",
    DBPass: "sdk-password",
    DBHost: "sdk-db",
    DBPort: 5432,
}

// ... or you can use the interface
// DB := database.NewDatabase(nil)
//
// func (db *Database) Config(dbname string, dbuser string, dbpass string, dbhost string, dbport int) *Database

// The database package exposes NewDatabase(config) *Database,
// which returns the database manager struct
// type Database struct {
// 	Engine    *gorm.DB
// 	DBadapter DBAdapter
// 	DBconfig  DBConfig
// 	LogLvl    logger.LogLevel
// 	log       zerolog.Logger
// }

DB := database.NewDatabase(config)

// Now we can choose the type of connection we need and initialize it.

// For Postgres...
database.DBManager = DB.Adapter(&database.PostgreAdapter{}).LogLevel(logger.Error).Init()

// For SQLite...
database.DBManager = DB.Adapter(&database.SQLiteAdapter{}).LogLevel(logger.Error).Init()

// For MySQL...
database.DBManager = DB.Adapter(&database.MySQLAdapter{}).LogLevel(logger.Error).Init()
``` 

The `database.DBManager` is a global variable declared in the database package. It stores the pointer to a Database struct containing an active DB connection.

Once the connection is established, developers can perform migration and seed operations on the connected database via the interfaces provided by the package.

Run a migration


```
// The database package exposes the NewMigration(*database.Database) *database.Migration
// interface that creates a new migration object needed to create a database structure
// based on the models provided by the developer.

// The Migration object expects a live database connection.
migration := database.NewMigration(database.DBManager)

// We can then create a map of all the models to migrate.
migrationModels := []interface{}{
	&models.DBSeeder{},
	&models.User{},
}

// Via the migrations AddModels(...interface{}) *database.Migration
// interface, you can add the model map to the migration object.
migration.AddModels(migrationModels...)

// Run the migration.
migration.Run()
```

### Run the seeders

Once the database has taken the intended shape, you can use the modules seed engine to populate some test data. In order to specify the seed logic, each model that will be seeded needs to inherit the SeedModel() interface, thus becoming a seeder model.

```
type User struct {
	ID   uint32
	Name string `json:"name"`
}

// Idea: what if we create interface for seeding data from models
// and then implement the seed logic in each model for that model

// Ok here's an example...

func (u *User) SeedModel() error {
	// seed logic here
	seeder := "seed_test_users"
	result := database.DBManager.Engine.Where("seeder_name", seeder).First(&DBSeeder{})
	if result.Error == gorm.ErrRecordNotFound {
		var users []*User
		user := &User{
			Name: "John Doe",
		}

		users = append(users, user)

		if result := database.DBManager.Engine.Create(&users); result.Error != nil {
			return result.Error
		}

		if result := database.DBManager.Engine.Create(&DBSeeder{SeederName: seeder}); result.Error != nil {
			return result.Error
		}
	}

	return nil
}
```

Instead of burying this logic deep inside the database package, we can use native Go interfaces to expose pieces of logic all the way up to the application layer. This way, developers can define the seed data for each model when needed.

```
// The database package exposes the NewSeeder(*database.Database) *database.Seeder
// interface, which creates a new seeder object for populating the database based on
// the models provided by the developer.
seeder := database.NewSeeder(database.DBManager)

// The AddSeeder(...) method can take one or more objects of DBSeeder type, and any model
// struct that implements the SeedModel() error interface is considered of DBSeeder type.
seeder.AddSeeder(&models.User{})

// Finally, the seeders are run.
seeder.Run()
```

### Data operations

Once our database is completely configured, we can start performing operations on its data. For that purpose, the module provides a generic repository interface that can add CRUD, filter, search, and pagination operations on any model by type embedding the generic CRUD repository type 

You can extend your custom repository. 

```
// Custom repository wrapper for User model
type UserRepository struct {
	dmp.Repository[*models.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Repository: *dmp.NewRepository(&models.User{}),
	}
}

All the data operations are now available in the UserRepository, so they can be called in a controller or a gRPC service implementation without any additional code. 

// REST
func (ctrl *UserController) Read(ctx *rest.Context) {
	ctx.SetContentType("application/json")
	ID := ctx.GetID()
	user := ctrl.User.Get(ID)

	if user == nil {
		ctx.JsonResponse(404, "User not found.")
		return
	}

	ctx.JsonResponse(200, user)
}

// gRPC
func (server *GRPCUserServer) GetUsers(ctx context.Context, req *connect.Request[gen.GetUsersRequest]) (response *connect.Response[gen.GetUsersResponse], err error) {
	users := []*gen.User{}

	usersData := server.User.GetAll()

	grpc.Map(usersData, &users)

	resp := &gen.GetUsersResponse{
		Users: users,
	}

	return connect.NewResponse(resp), nil
```

A developer can add custom wrappers and additional logic, in the repository implementation, add additional checks, format the data, etc.

```
func (repo *UserRepository) GetUsers() []*models.User {
	resp, err := repo.GetAll()
	if err == nil {
		return resp
	}
	return nil
}

```

The filter, pagination and search functionalities are implemented but not yet fully tested.