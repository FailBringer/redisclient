package redisClientAPI

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool          *redis.Pool
	redisServer   = "127.0.0.1:6379"
	redisPassword = ""
)

func main() {

}
func NewPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
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
}

func CreateNewClient(id, key, name, totalBudget string) string {
	if id == "" {
		return "Input id please"
	}
	if key == "" {
		return "Input key please"
	}
	if name == "" {
		return "Input name please"
	}
	if totalBudget == "" {
		return "Input budget please"
	}
	conn := GetConnection()
	defer conn.Close()
	check, err := redis.String(conn.Do("HGET", "client:"+id, "id"))
	if err != nil {
		if err != redis.ErrNil {
			fmt.Println(err)
		}
	}
	if check == "" {
		conn.Send("HSET", "client:"+id, "id", id)
		conn.Send("HSET", "client:"+id, "key", key)
		conn.Send("HSET", "client:"+id, "name", name)
		conn.Send("HSET", "client:"+id, "totalBudget", totalBudget)
		conn.Flush()
		return "Client created"
	} else {
		return "Id alredy exist"
	}
}

func DeleteClient(id string) string {
	conn := GetConnection()
	defer conn.Close()
	result, err := redis.Int(conn.Do("DEL", "client:"+id))
	if err != nil {
		fmt.Println(err)
	}
	if result == 1 {
		return "OK"
	} else {
		return "Something Wrong"
	}
}

func GetConnection() redis.Conn {
	pool := NewPool(redisServer, redisPassword)
	conn := pool.Get()
	return conn
}

func ClientAuth(id, key string) bool {
	conn := GetConnection()
	defer conn.Close()
	pass, err := redis.String(conn.Do("HGET", "client:"+id, "key"))
	if err != nil {
		fmt.Println(err)
	}
	if key == pass {
		return true
	} else {
		return false
	}
}

func GetTotalBudget(id, key string) string {
	if ClientAuth(id, key) {
		conn := GetConnection()
		defer conn.Close()
		totalBudget, err := redis.String(conn.Do("HGET", "client:"+id, "totalBudget"))
		if err != nil {
			fmt.Println(err)
			return err.Error()
		}
		return totalBudget
	} else {
		return "Access denied"
	}

}

func SetTotalBudget(id, key, budget string) string {
	if ClientAuth(id, key) {
		conn := GetConnection()
		defer conn.Close()
		_, err := redis.Int(conn.Do("HSET", "client:"+id, "totalBudget", budget))
		if err != nil {
			fmt.Println(err)
			return err.Error()
		}
		return "OK"
	} else {
		return "Access denied"
	}
}
