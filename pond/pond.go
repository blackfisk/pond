package pond

import (
	"io/ioutil"
	"log"
	"os"
	"bufio"
	"bytes"
	"net/http"
	"encoding/json"
	"runtime"

	"github.com/garyburd/redigo/redis"
)

type Pond struct {
	queue chan *Rock
}

func NewPond() *Pond {
	p := new(Pond)
	p.queue = make(chan *Rock)

	p.startWorkers()
	p.startBroadcasters()

	// Wake up conn!
	conn := pool.Get()
	conn.Do("PING")
	defer conn.Close()

	return p
}

func (p *Pond) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	body := req.Body

        switch req.Method {
        case "POST":
                if body != nil {
                        bytes, _ := ioutil.ReadAll(body)
                        rock := NewRock(bytes)

                        if !rock.alreadySent() {
                                p.queue <- rock
                        }
                }

        case "GET":
                conn := pool.Get()
                defer conn.Close()

                keys, _ := redis.Strings(conn.Do("KEYS", message_key + ":*"))
                messages := make([]string, 0)

                for _, key := range keys {
                        msg, _ := redis.String(conn.Do("GET", key))
                        messages = append(messages, msg)
                }

                output, _ := json.Marshal(messages)

                w.Write(output)
        }
}

func (p *Pond) storeMessage(msg []byte) {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", queue_key, msg)
	if err != nil {
		panic(err)
	}

	log.Println("---> Message stored!")
}

func (p *Pond) worker(i int) {
	log.Println("---> Starting the worker", i)

	var rock *Rock

	for {
		rock = <-p.queue
		if !rock.alreadySent() {
			log.Printf("----> [w:%d] Incoming message!", i)
			p.storeMessage(rock.Message)
		}
	}
}

func (p *Pond) broadcaster(i int) {
	log.Println("---> Starting the broadcaster", i)
	conn := pool.Get()
	defer conn.Close()

	for {
		msg, _ := redis.Bytes(conn.Do("BRPOPLPUSH", queue_key, backup_key, 5))
		if msg != nil {
			rock := NewRock(msg)

			if !rock.alreadySent() {
				go p.sendToTheRiver(rock)
				rock.StoreForReading()
				rock.FlagAsSent()

				log.Printf("----> [b:%d] Sent message", i)
			}
		}

	}
}

func (p *Pond) sendToTheRiver(rock *Rock) {
        friendsFile, _ := os.Open(".ponds")
        scanner := bufio.NewScanner(friendsFile)
        message := bytes.NewReader(rock.Message)

        for scanner.Scan() {
                friend := scanner.Text()
                http.Post("http://" + friend, "text/plain", message)
        }
}

func (p *Pond) startWorkers() {
	cpu_count := runtime.NumCPU() / 2
	runtime.GOMAXPROCS(cpu_count)

	log.Printf("--> Starting %d workers\n", cpu_count)

	for i := 0; i < cpu_count; i++ {
		go p.worker(i)
	}
}

func (p *Pond) startBroadcasters() {
	cpu_count := runtime.NumCPU() / 2
	log.Printf("--> Starting %d broadcasters\n", cpu_count)

	for i := 0; i < cpu_count; i++ {
		go p.broadcaster(i)
	}
}
