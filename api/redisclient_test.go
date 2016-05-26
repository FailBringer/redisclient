package redisClientAPI

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedisConnection(t *testing.T) {
	conn := GetConnection()
	defer conn.Close()
	assert.Nil(t, conn.Err())
}

func TestClientCreation(t *testing.T) {
	conn := GetConnection()
	defer conn.Close()
	conn.Send("HGET", "client:test", "name")
	conn.Flush()
	status, err := redis.String(conn.Receive())
	if err != nil {
		fmt.Println(err)
	}
	if status != "" {
		reply := DelTestData(conn)
		assert.Equal(t, reply, 1, "TestClientCreation DelTestData assert")
	}
	CreateTestData()
	conn.Send("HGET", "client:test", "name")
	conn.Flush()
	result, err := redis.String(conn.Receive())
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, result, "Andrew", "TestClientCreation result assert")
	DelTestData(conn)
}

func TestCreateNewClient(t *testing.T) {
	status := CreateNewClient("test", "pass", "Andrew", "2890")
	assert.Equal(t, status, "Client created")
	DeleteClient("test")
}

func TestCreateNewClientWithoutSomeData(t *testing.T) {
	status := CreateNewClient("", "pass", "Andrew", "7568")
	assert.Equal(t, status, "Input id please", "assert without id")
	status = CreateNewClient("test", "", "Andrew", "7568")
	assert.Equal(t, status, "Input key please", "assert without key")
	status = CreateNewClient("test", "pass", "", "7568")
	assert.Equal(t, status, "Input name please", "assert without name")
	status = CreateNewClient("test", "pass", "Andrew", "")
	assert.Equal(t, status, "Input budget please", "assert without budget")
}

func TestCreateNewClientWhenIdAlredyExist(t *testing.T) {
	CreateTestData()
	status := CreateNewClient("test", "passKey", "Max", "78965")
	assert.Equal(t, status, "Id alredy exist")
	DelTestData(GetConnection())
}

func DelTestData(conn redis.Conn) int {
	reply, err := redis.Int(conn.Do("DEL", "client:test"))
	if err != nil {
		fmt.Println(err)
	}
	return reply
}

func CreateTestData() {
	CreateNewClient("test", "pass", "Andrew", "750000")
}

func TestClientDeleting(t *testing.T) {
	CreateTestData()
	result := DeleteClient("test")
	assert.Equal(t, result, "OK", "TestClientDeleting DeleteClient result assert")
	conn := GetConnection()
	defer conn.Close()
	reply, err := redis.String(conn.Do("HGET", "client:test", "name"))
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, reply, "", "TestClientDeleting reply assert")
}

func TestClientAuth(t *testing.T) {
	CreateTestData()
	result := ClientAuth("test", "pass")
	assert.Equal(t, result, true, "TestClientAuth result assert")
	DelTestData(GetConnection())
}

func TestGetTotalButget(t *testing.T) {
	CreateTestData()
	result := GetTotalBudget("test", "pass")
	assert.Equal(t, result, "750000", "TestGetTotalButget result assert")
	conn := GetConnection()
	defer conn.Close()
	reply, err := redis.String(conn.Do("HGET", "client:test", "totalBudget"))
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, reply, "750000", "TestGetTotalButget reply assert")
	DelTestData(GetConnection())
}

func TestSetTotalBudget(t *testing.T) {
	CreateTestData()
	result := SetTotalBudget("test", "pass", "745000")
	assert.Equal(t, result, "OK", "TestSetTotalBudget result assert")
	result = GetTotalBudget("test", "pass")
	assert.Equal(t, result, "745000", "TestSetTotalBudget GetTotalBudget assert")
	DelTestData(GetConnection())
}
