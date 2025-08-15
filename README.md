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

Create a docker network
```
docker network create cassandra
```
To inspect a docker network
```
docker network inspect cassandra
```
Run the Cassandra docker container
```
docker run --name cassandra -p 9042:9042 -d cassandra:4.1
```

To manually connect the cassandra container to the cassandra network
```
docker network connect cassandra cassandra
```
Run this to run the `data.cql` script to seed some data
```
docker run --rm --network cassandra -v "$(pwd)/scripts/data.cql:/scripts/data.cql" cassandra cqlsh cassandra 9042 -f /scripts/data.cql
```

Run this to get a CQL shell to query the database, you may need to wait a little while for the database to start. Use `ctrl + D` to exit the CQL shell.
```
docker run --rm -it --network cassandra cassandra:latest cqlsh cassandra 9042 --cqlversion='3.4.6'
```
Example select statement
```
SELECT * FROM filmclub.users_by_id;
```