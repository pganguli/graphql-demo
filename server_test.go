package main

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"encoding/json"
	"net/http"

	"github.com/pganguli/hnews/pkg/jwt"
)

type GQL_Query struct {
	OperationName string
	Query         string
}

type GQL_Response struct {
	Data   map[string]interface{}
	Errors []map[string]interface{}
}

type httpHeaders map[string]string

var GQL_Url = "http://localhost:8080/graphql"

func GqlQuery(gql_query GQL_Query, headers httpHeaders) (response GQL_Response) {
	json_str, err := json.Marshal(gql_query)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, GQL_Url, bytes.NewReader(json_str))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	resp_body := new(bytes.Buffer)
	if _, err := resp_body.ReadFrom(resp.Body); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(resp_body.Bytes(), &response); err != nil {
		panic(err)
	}

	return response
}

func CreateUser(username string, password string) (response GQL_Response) {
	operationName := "create_user"
	query := fmt.Sprintf("mutation %v{createUser(input:{username:\"%v\",password:\"%v\"})}",
		operationName, username, password)

	data := GQL_Query{
		OperationName: operationName,
		Query:         query,
	}

	return GqlQuery(data, nil)
}

func CreateLink(title string, address string, bearerToken httpHeaders) (response GQL_Response) {
	operationName := "create_link"
	query := fmt.Sprintf("mutation %v{createLink(input:{title:\"%v\",address:\"%v\"}){id title address user{id username}}}",
		operationName, title, address)

	data := GQL_Query{
		OperationName: operationName,
		Query:         query,
	}

	return GqlQuery(data, bearerToken)
}

func ListLinks() (response GQL_Response) {
	operationName := "list_links"
	query := fmt.Sprintf("query %v{links{id title address user{id username}}}",
		operationName)

	data := GQL_Query{
		OperationName: operationName,
		Query:         query,
	}

	return GqlQuery(data, nil)
}

func Login(username string, password string) (response GQL_Response) {
	operationName := "login"
	query := fmt.Sprintf("mutation %v{login(input:{username:\"%v\",password:\"%v\"})}",
		operationName, username, password)

	data := GQL_Query{
		OperationName: operationName,
		Query:         query,
	}

	return GqlQuery(data, nil)
}

func RefreshToken(token string) (response GQL_Response) {
	operationName := "refresh_token"
	query := fmt.Sprintf("mutation %v{refreshToken(input:{token:\"%v\"})}",
		operationName, token)

	data := GQL_Query{
		OperationName: operationName,
		Query:         query,
	}

	return GqlQuery(data, nil)
}

func TestCreateUser(t *testing.T) {
	response := CreateUser("user0", "pa$$word")
	if response.Errors != nil {
		t.Errorf("%+v", response.Errors)
	}

	authToken, ok := response.Data["createUser"].(string)
	if !ok {
		t.Errorf("Invalid data:\n%+v", response.Data)
	}

	userName, err := jwt.ParseToken(authToken)
	if err != nil {
		t.Error(err)
	}
	if userName != "user0" {
		t.Errorf("Claimed user name is incorrect: %q", userName)
	}

	duplicate_user_text := "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)"

	response = CreateUser("user0", "pa$$word")
	if response.Errors == nil {
		t.Errorf("Did not get expected error: %q", duplicate_user_text)
	} else {
		message := response.Errors[0]["message"].(string)
		if message != duplicate_user_text {
			t.Errorf("Expected error: %q\nGot error: %q", duplicate_user_text, message)
		}
	}
}

func TestCreateUserLink(t *testing.T) {
	response := CreateUser("user1", "pa$$word")
	if response.Errors != nil {
		t.Fatalf("%+v", response.Errors)
	}

	authToken := response.Data["createUser"].(string)

	access_denied_text := "access denied"

	// No bearer token
	response = CreateLink("title1", "ww1.example.net", httpHeaders{})
	if response.Errors == nil {
		t.Fatalf("Did not get expected error: %q", access_denied_text)
	} else {
		message := response.Errors[0]["message"].(string)
		if message != access_denied_text {
			t.Errorf("Expected error: %q\nGot error: %q", access_denied_text, message)
		}
	}

	// Invalid bearer token
	response = CreateLink("title1", "ww1.example.net", httpHeaders{"Authorization": "bearer " + authToken})
	if response.Errors == nil {
		t.Fatalf("Did not get expected error: %q", access_denied_text)
	} else {
		message := response.Errors[0]["message"].(string)
		if message != access_denied_text {
			t.Errorf("Expected error: %q\nGot error: %q", access_denied_text, message)
		}
	}

	// Valid bearer token
	response = CreateLink("title1", "ww1.example.net", httpHeaders{"Authorization": "Bearer " + authToken})
	if response.Errors != nil {
		t.Fatalf("%+v", response.Errors)
	}

	data, ok := response.Data["createLink"].(map[string]interface{})
	if !ok {
		t.Errorf("Invalid data:\n%+v", data)
	}

	address, ok := data["address"].(string)
	if !ok || address != "ww1.example.net" {
		t.Errorf("Invalid address: %q", address)
	}

	id_s, ok := data["id"].(string)
	if _, err := strconv.Atoi(id_s); !ok || err != nil {
		t.Errorf("Invalid id: %q", id_s)
	}

	title, ok := data["title"].(string)
	if !ok || title != "title1" {
		t.Errorf("Invalid title: %q", title)
	}

	user := data["user"].(map[string]interface{})

	user_id_s, ok := user["id"].(string)
	if _, err := strconv.Atoi(user_id_s); !ok || err != nil {
		t.Errorf("Invalid user_id: %q", user_id_s)
	}

	user_username, ok := user["username"].(string)
	if !ok || user_username != "user1" {
		t.Errorf("Invalid user_username: %q", user_username)
	}
}

func TestListLinks(t *testing.T) {
	response := CreateUser("user2", "pa$$word")
	if response.Errors != nil {
		t.Fatalf("%+v", response.Errors)
	}

	authToken := response.Data["createUser"].(string)
	bearerToken := httpHeaders{"Authorization": "Bearer " + authToken}

	for i := 1; i <= 5; i++ {
		response = CreateLink(
			"title"+strconv.Itoa(i),
			"ww"+strconv.Itoa(i)+".example.net",
			bearerToken)
		if response.Errors != nil {
			t.Fatalf("%+v", response.Errors)
		}
	}

	response = ListLinks()
	if response.Errors != nil {
		t.Fatalf("%+v", response.Errors)
	}

	links, ok := response.Data["links"].([]interface{})
	if !ok {
		t.Errorf("Invalid links:\n%+v", response.Data)
	}

	count := 0
	for _, link := range links {
		user, ok := link.(map[string]interface{})["user"]
		if !ok {
			t.Errorf("Invalid user:\n%+v", link)
		}

		username, ok := user.(map[string]interface{})["username"]
		if !ok {
			t.Errorf("Invalid username:\n%+v", user)
		}

		if username == "user2" {
			count++
		}
	}

	if count != 5 {
		t.Errorf("Expected 5 links; got: %d", count)
	}
}

func TestLogin(t *testing.T) {
	response := CreateUser("user3", "pa$$word")
	if response.Errors != nil {
		t.Fatalf("%+v", response.Errors)
	}

	authToken := response.Data["createUser"].(string)

	response = Login("user3", "pa$$word")
	if response.Errors != nil {
		t.Fatalf("%+v", response.Errors)
	} else {
		loginToken := response.Data["login"].(string)
		if loginToken != authToken {
			t.Errorf("Expected token: %q\nGot token: %q", authToken, loginToken)
		}
	}
}

func TestRefreshToken(t *testing.T) {
	response := CreateUser("user4", "pa$$word")
	if response.Errors != nil {
		t.Fatalf("%+v", response.Errors)
	}

	authToken := response.Data["createUser"].(string)

	response = RefreshToken(authToken)
	if response.Errors != nil {
		t.Fatalf("%+v", response.Errors)
	} else {
		refreshedToken := response.Data["refreshToken"].(string)
		if refreshedToken != authToken {
			t.Errorf("Expected token: %q\nGot token: %q", authToken, refreshedToken)
		}
	}
}
