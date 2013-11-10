package pond

import (
        "fmt"
	"crypto/sha1"
	"encoding/base64"

	"github.com/garyburd/redigo/redis"
)

type Rock struct {
	Message []byte
	Hash    string
}

func NewRock(msg []byte) *Rock {
	r := new(Rock)
	r.Message = msg
	r.Hash = r.MessageHash(string(msg))

	return r
}

func (r *Rock) StoreForReading() {
	conn := pool.Get()
	defer conn.Close()

        key := fmt.Sprintf("%s:%s", message_key, r.Hash)

        conn.Do("SETEX", key, 3600*24*2, r.Message)
}

func (r *Rock) MessageHash(msg string) string {
	h := sha1.New()
	h.Write([]byte(msg))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return sha
}

func (r *Rock) FlagAsSent() {
	conn := pool.Get()
	defer conn.Close()

	conn.Do("SADD", sent_key, r.Hash)
	conn.Do("LPOP", backup_key)
}

func (r *Rock) alreadySent() bool {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", sent_key))
	if err != nil {
		panic(err)
	}

	if !exists {
		return false
	}

	hash := string(r.Hash)
	already_sent, _ := redis.Bool(conn.Do("SISMEMBER", sent_key, hash))

	return already_sent
}
