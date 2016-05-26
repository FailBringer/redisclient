package main

import (
	"net/http"
	"redisclient/api"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	action := GetAction(r)
	switch action {
	default:
		{
			w.Write([]byte("HELLO"))
		}

	case "GetTotalBudget":
		{
			id := GetId(r)
			key := GetKey(r)
			totalbudget := GetTotalBudget(id, key)
			w.Write([]byte(totalbudget))
		}

	case "SetTotalBudget":
		{
			id := GetId(r)
			key := GetKey(r)
			budget := GetBudget(r)
			status := SetTotalBudget(id, key, budget)
			w.Write([]byte(status))
		}

	case "CreateNewClient":
		{
			id := GetId(r)
			key := GetKey(r)
			name := GetName(r)
			budget := GetBudget(r)
			status := CreateNewClient(id, key, name, budget)
			w.Write([]byte(status))
		}
	}
}

func main() {
	http.HandleFunc("/client", Handler)
	http.ListenAndServe(":8080", nil)
}

func GetAction(r *http.Request) string {
	r.ParseForm()
	action := r.Form.Get("action")
	return action
}

func GetId(r *http.Request) string {
	r.ParseForm()
	id := r.Form.Get("id")
	return id
}

func GetName(r *http.Request) string {
	r.ParseForm()
	name := r.Form.Get("name")
	return name
}

func GetBudget(r *http.Request) string {
	r.ParseForm()
	budget := r.Form.Get("budget")
	return budget
}

func GetKey(r *http.Request) string {
	r.ParseForm()
	key := r.Form.Get("key")
	return key
}

func GetTotalBudget(id, key string) string {
	totalbudget := redisClientAPI.GetTotalBudget(id, key)
	return totalbudget
}

func SetTotalBudget(id, key, budget string) string {
	status := redisClientAPI.SetTotalBudget(id, key, budget)
	return status
}

func CreateNewClient(id, key, name, budget string) string {
	status := redisClientAPI.CreateNewClient(id, key, name, budget)
	return status
}
