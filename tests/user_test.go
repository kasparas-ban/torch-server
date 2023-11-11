package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"torch/torch-server/models"
	"torch/torch-server/testutil"

	"github.com/stretchr/testify/assert"
)

var numUsersInDB = 3

func TestAddUser(t *testing.T) {
	testutil.CleanAllTables()
	testutil.SeedDB()

	// Router setup
	userID = uint64(1)
	w, c, router := RouterSetup(userID)

	// Adding a new user
	requestJson := []byte(`
		{
			"username": "new_user",
			"email": "test_email@gmail.com",
			"birthday": "2000-01-01",
			"gender": "MALE",
			"countryID": 130,
			"city": "Vilnius",
			"description": "Some description about me."
		}
	`)

	var requestBody models.NewUser
	err := json.Unmarshal(requestJson, &requestBody)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-user", bytes.NewReader(requestJson))
	router.ServeHTTP(w, c.Request)

	var addedUser models.ExistingUser
	if err := json.Unmarshal(w.Body.Bytes(), &addedUser); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, requestBody.Username, addedUser.Username)
	assert.Equal(t, requestBody.Email, addedUser.Email)
	assert.Equal(t, requestBody.Birthday, addedUser.Birthday)
	assert.Equal(t, requestBody.Gender, addedUser.Gender)
	assert.NotEmpty(t, addedUser.CountryCode)
	assert.Equal(t, requestBody.City, addedUser.City)
	assert.Equal(t, requestBody.Description, addedUser.Description)

	// Router setup
	userID = uint64(numUsersInDB + 1)
	w, c, router = RouterSetup(userID)

	// Getting user info
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user-info", nil)
	router.ServeHTTP(w, c.Request)

	var user models.ExistingUser
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, requestBody.Username, user.Username)
	assert.Equal(t, requestBody.Email, user.Email)
	assert.Equal(t, requestBody.Birthday, user.Birthday)
	assert.Equal(t, requestBody.Gender, user.Gender)
	assert.NotEmpty(t, user.CountryCode.Val)
	assert.Equal(t, requestBody.City, user.City)
	assert.Equal(t, requestBody.Description, user.Description)

	testutil.CleanAllTables()
}

func TestDeleteUser(t *testing.T) {
	testutil.CleanAllTables()
	testutil.SeedDB()

	// Router setup
	userID = uint64(0)
	w, c, router := RouterSetup(userID)

	// Adding a new user
	requestJson := []byte(`
		{
			"username": "new_user",
			"email": "test_email@gmail.com",
			"birthday": "2000-01-01",
			"gender": "MALE",
			"countryID": 130,
			"city": "Vilnius",
			"description": "Some description about me."
		}
	`)

	var requestBody models.NewUser
	err := json.Unmarshal(requestJson, &requestBody)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-user", bytes.NewReader(requestJson))
	router.ServeHTTP(w, c.Request)

	var addedUser models.ExistingUser
	if err := json.Unmarshal(w.Body.Bytes(), &addedUser); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	// Router setup
	userID = uint64(numUsersInDB + 1)
	w, c, router = RouterSetup(userID)

	// Deleting the user
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/delete-user", nil)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)

	// Router setup
	userID = uint64(numUsersInDB + 1)
	w, c, router = RouterSetup(userID)

	// Getting deleted user info
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user-info", nil)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	testutil.CleanAllTables()
}

func TestUpdateUser(t *testing.T) {
	testutil.CleanAllTables()
	testutil.SeedDB()

	// Router setup
	userID = uint64(0)
	w, c, router := RouterSetup(userID)

	// Adding a new user
	requestJson := []byte(`
		{
			"username": "new_user",
			"email": "test_email@gmail.com",
			"birthday": "2000-01-01",
			"gender": "MALE",
			"countryID": 130,
			"city": "Vilnius",
			"description": "Some description about me."
		}
	`)

	var requestBody models.NewUser
	err := json.Unmarshal(requestJson, &requestBody)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-user", bytes.NewReader(requestJson))
	router.ServeHTTP(w, c.Request)

	var addedUser models.ExistingUser
	if err := json.Unmarshal(w.Body.Bytes(), &addedUser); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	// Router setup
	userID = uint64(numUsersInDB + 1)
	w, c, router = RouterSetup(userID)

	// Updating the user
	requestJson = []byte(`
		{
			"username": "updated_new_user",
			"email": "updated_test_email@gmail.com",
			"birthday": "2010-01-01",
			"gender": "FEMALE",
			"countryID": 129,
			"city": "Kaunas",
			"description": "Updated description."
		}
	`)

	var updatedBody models.NewUser
	err = json.Unmarshal(requestJson, &updatedBody)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = httptest.NewRequest(http.MethodPut, "/api/update-user", bytes.NewReader(requestJson))
	router.ServeHTTP(w, c.Request)

	// Router setup
	userID = uint64(numUsersInDB + 1)
	w, c, router = RouterSetup(userID)

	// Getting the updated user's info
	c.Request = httptest.NewRequest(http.MethodGet, "/api/user-info", nil)
	router.ServeHTTP(w, c.Request)

	var updatedUser models.ExistingUser
	if err := json.Unmarshal(w.Body.Bytes(), &updatedUser); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, updatedBody.Username, updatedUser.Username)
	assert.Equal(t, updatedBody.Email, updatedUser.Email)
	assert.Equal(t, updatedBody.Birthday, updatedUser.Birthday)
	assert.Equal(t, updatedBody.Gender, updatedUser.Gender)
	assert.Equal(t, updatedBody.City, updatedUser.City)
	assert.Equal(t, updatedBody.Description, updatedUser.Description)

	testutil.CleanAllTables()
}
