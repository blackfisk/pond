package pond

import (
	"os"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	server = os.Getenv("OPENREDIS_URL")[8:]
	parts  = strings.Split(server, "@")

	password  = parts[0][1:]
	redis_url = parts[1]

        messages_key = "pond:messages"
        queue_key  = "pond:queue"
        backup_key = "pond:backup"
        sent_key   = "pond:sent"
)

var pool = &redis.Pool{
	IdleTimeout: 240 * time.Second,
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", redis_url)
		if err != nil {
			return nil, err
		}

		if len(password) > 0 {
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
		}
		return c, err
	},

	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	},
}
