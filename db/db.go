package db

import (
	"github.com/gocql/gocql"
	"log"
	"os"
	"strings"
	"time"
)

var session *gocql.Session

type ChatMessage struct {
	UserID          string
	UserName        string
	UserDisplayName string
	Message         string
	Channel         string
	Time            time.Time
	Bits            int
}

func Connect() {
	// read hosts from environment
	hostsenv := os.Getenv("CASSANDRA_HOSTS")
	hostsenv = strings.Replace(hostsenv, " ", "", -1)
	hosts := strings.Split(hostsenv, ",")

	log.Printf("Connecting to Cassandra hosts %v\n", hosts)

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "twitch"
	cluster.Consistency = gocql.Quorum
	s, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	session = s

	// Create message table (if not existing)
	createTableQuery := session.Query("CREATE TABLE IF NOT EXISTS message (id uuid PRIMARY KEY, user_id text, user_name text, user_display_name text, message text, channel text, time timestamp, bits int)")
	err := createTableQuery.Exec()
	if err != nil {
		panic(err)
	}
}

func Disconnect() {
	if session != nil {
		session.Close()
	}
}

func InsertMessage(message ChatMessage) error {
	queryStr := `INSERT INTO message (id, user_id, user_name, user_display_name, message, channel, time, bits) VALUES (uuid(), ?, ?, ?, ?, ?, ?, ?)`
	query := session.Query(queryStr,
		message.UserID,
		message.UserName,
		message.UserDisplayName,
		message.Message,
		message.Channel,
		message.Time,
		message.Bits)
	return query.Exec()
}
