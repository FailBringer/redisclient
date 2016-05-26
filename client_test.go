package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"redisclient/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAction(t *testing.T) {
	req := CreateTestRequest("action=test")
	action := GetAction(req)
	assert.Equal(t, action, "test", "TestGetAction assert")
}

func TestGetId(t *testing.T) {
	req := CreateTestRequest("id=test")
	id := GetId(req)
	assert.Equal(t, id, "test", "TestGetId assert")
}

func TestGetKey(t *testing.T) {
	req := CreateTestRequest("key=pass")
	key := GetKey(req)
	assert.Equal(t, key, "pass", "TestGetId assert")
}

func TestHandler(t *testing.T) {
	req := CreateTestRequest()
	w := httptest.NewRecorder()
	Handler(w, req)
	assert.Equal(t, w.Body.String(), "HELLO", "TestHandler assert")
}

func TestHandlerWithWrongFields(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "5000")
	req := CreateTestRequest("action=GetTotalBudget", "key=pass")
	w := httptest.NewRecorder()
	Handler(w, req)
	assert.Equal(t, w.Body.String(), "Access denied", "TestHandlerWithWrongFields assert")
	redisClientAPI.DeleteClient("test")
}

func TestHandlerGetTotalBudget(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "5000")
	req := CreateTestRequest("action=GetTotalBudget", "id=test", "key=pass")
	w := httptest.NewRecorder()
	Handler(w, req)
	assert.Equal(t, w.Body.String(), "5000", "TestHandlerGetTotalBudget assert")
	redisClientAPI.DeleteClient("test")
}

func TestHandlerGetTotalBudgetWithWrongKey(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "5000")
	req := CreateTestRequest("action=GetTotalBudget", "id=test", "key=wrongPass")
	w := httptest.NewRecorder()
	Handler(w, req)
	assert.Equal(t, w.Body.String(), "Access denied", "TestHandlerGetTotalBudgetWithWrongKey assert")
	redisClientAPI.DeleteClient("test")
}

func TestHandlerSetTotalBudget(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "5000")
	req := CreateTestRequest("action=SetTotalBudget", "id=test", "key=pass", "budget=3000")
	w := httptest.NewRecorder()
	Handler(w, req)
	assert.Equal(t, w.Body.String(), "OK", "TestHandlerSetTotalBudget assert")
	redisClientAPI.DeleteClient("test")
}

func TestHandlerSetTotalBudgetWithWrongKey(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "5000")
	req := CreateTestRequest("action=SetTotalBudget", "id=test", "key=wrongPass", "budget=100")
	w := httptest.NewRecorder()
	Handler(w, req)
	assert.Equal(t, w.Body.String(), "Access denied", "TestHandlerSetTotalBudgetWithWrongKey assert")
	redisClientAPI.DeleteClient("test")
}

func TestHandlerCreateNewClient(t *testing.T) {
	req := CreateTestRequest("action=CreateNewClient", "id=test", "key=pass", "name=test", "budget=1000")
	w := httptest.NewRecorder()
	Handler(w, req)
	assert.Equal(t, w.Body.String(), "Client created", "TestHandlerCreateNewClient assert")
	redisClientAPI.DeleteClient("test")

}

func TestHandlerCreateNewClientWithError(t *testing.T) {
	req := CreateTestRequest("action=CreateNewClient", "key=pass", "name=test", "budget=1000")
	w := httptest.NewRecorder()
	Handler(w, req)
	assert.Equal(t, w.Body.String(), "Input id please", "TestHandlerCreateNewClientWithError assert")
	redisClientAPI.DeleteClient("test")

}

func CreateTestRequest(args ...string) *http.Request {
	reqString := "http://127.0.0.1:8080/client?"
	for _, s := range args {
		reqString += s
		reqString += "&"
	}
	req, err := http.NewRequest("GET", reqString, nil)
	if err != nil {
		fmt.Println(err)
	}
	return req
}

func TestGetTotalBudget(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "70000")
	totalbudget := GetTotalBudget("test", "pass")
	assert.Equal(t, totalbudget, "70000", "TestGetTotalBudget assert")
	redisClientAPI.DeleteClient("test")
}

func TestGetTotalBudgetWhithWrongKey(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "70000")
	totalbudget := GetTotalBudget("test", "wrongPass")
	assert.Equal(t, totalbudget, "Access denied", "TestGetTotalBudgetWhithWrongKey assert")
	redisClientAPI.DeleteClient("test")
}

func TestGetTotalBudgetWhithWrongId(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "70000")
	totalbudget := GetTotalBudget("wrongId", "pass")
	assert.Equal(t, totalbudget, "Access denied", "TestGetTotalBudgetWhithWrongId assert")
	redisClientAPI.DeleteClient("test")
}

func TestSetTotalBudget(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "70000")
	status := SetTotalBudget("test", "pass", "6000")
	assert.Equal(t, status, "OK", "TestSetTotalBudget assert")
	redisClientAPI.DeleteClient("test")
}

func TestSetTotalBudgetWithWrongKey(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "70000")
	status := SetTotalBudget("test", "wrongKey", "6000")
	assert.Equal(t, status, "Access denied", "TestSetTotalBudgetWithWrongKey assert")
	redisClientAPI.DeleteClient("test")
}

func TestSetTotalBudgetWithWrongId(t *testing.T) {
	redisClientAPI.CreateNewClient("test", "pass", "Andrew", "70000")
	status := SetTotalBudget("wrongId", "pass", "6000")
	assert.Equal(t, status, "Access denied", "TestSetTotalBudgetWithWrongId assert")
	redisClientAPI.DeleteClient("test")
}
