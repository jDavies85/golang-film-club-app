## Notes
In lieu of a proper readme these are just a collection of notes to help me with development.

### Development environment
TODO:
- prerequisite installs
- recommended VS code extensions etc.
### Running the api
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
To run the api enter the following command
```
go run .
```
### Cassandra

Create a network
```
docker network create cassandra
```

Run the Cassandra docker container
```
docker run --rm -d --name cassandra --hostname cassandra --network cassandra cassandra
```
Run this to run the `data.cql` script
```
docker run --rm --network cassandra -v "$(pwd)/data.cql:/scripts/data.cql" cassandra:latest cqlsh cassandra 9042 -f /scripts/data.cql
```

Run this to get a CQL shell to query the database
```
docker run --rm -it --network cassandra cassandra:latest cqlsh cassandra 9042 --cqlversion='3.4.7'
```
Example select statement
```
SELECT * FROM store.shopping_cart;
```