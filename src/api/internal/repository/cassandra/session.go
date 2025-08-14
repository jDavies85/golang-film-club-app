package cassandra

import (
	"time"

	"github.com/gocql/gocql"
)

func NewSession(hosts []string, keyspace, consistency string, timeoutSec int) (*gocql.Session, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Timeout = time.Duration(timeoutSec) * time.Second
	switch consistency {
	case "LOCAL_QUORUM":
		cluster.Consistency = gocql.LocalQuorum
	default:
		cluster.Consistency = gocql.Quorum
	}
	return cluster.CreateSession()
}
