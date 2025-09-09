## Notes
In lieu of a proper readme these are just a collection of notes to help me with development.

### Development environment
TODO:
- prerequisite installs
- recommended VS code extensions etc.

### Running the api
To run the api enter the following command
```
cd src/api
go run ./cmd/server
```
### Go development tips
To add a third party dependency, first `go get` it.
```
go get "github.com/gin-gonic/gin"
```
Then add it to the imports at the top of the file:
```
import (
	"net/http"

	"github.com/gin-gonic/gin"
)
```

TODO talk about devauth setup
### Cassandra

Get the latest version of Cassandra
```
docker pull cassandra:latest
```

Run the container
```
docker run --name cassandra -p 9042:9042 -d cassandra:latest
```

Run this to get a CQL shell to query the database, you may need to wait a little while for the database to start. Use `ctrl + D` to exit the CQL shell.
```
docker exec -it cassandra cqlsh
```

Run this to run the `data.cql` script to seed some data
```
cd src\api
Get-Content .\scripts\data.cql | docker exec -i cassandra cqlsh
```
Then star a cql session and run the following query
```
SELECT table_name 
FROM system_schema.tables 
WHERE keyspace_name = 'filmclub';
```
and you should see the following output
```
 table_name
----------------------
 club_members_by_club
     film_clubs_by_id
    membership_guards
   user_clubs_by_user
        users_by_auth
          users_by_id
```

Example select statement
```
SELECT * FROM filmclub.users_by_id;
```