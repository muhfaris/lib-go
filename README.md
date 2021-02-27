# lib-go

lib-go is a Golang library for saving Your time.


## Usage
### PostgreSQL

```
import (
    ...
    ...
    "github.com/muhfaris/lib-go/psql"
)

func initializePSQL() *sql.DB {
	dbOptions := psql.DBOptions{
		Host:     app.Database.Host,
		Port:     app.Database.Port,
		Username: app.Database.Username,
		Password: app.Database.Password,
		DBName:   app.Database.Name,
		SSLMode:  app.Database.SSLMode,
	}

	conn, err := psql.Connect(&dbOptions)
	if err != nil {
		log.Fatalln("Database:", err)
	}

   return conn
}

```

#### ParseBody
```
import (
    ...
    ...
	"github.com/pixelhousestudio/library-microservice/utils"
)

type User struct {
    Username string
    Password string
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user User

	err := utils.ParseBody(ctx, r, &user)
	if err != nil {
        # error
	}

    # TODO something
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT license](https://github.com/muhfaris/lib-go/blob/master/LICENSE).
