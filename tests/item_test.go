package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"torch/torch-server/controllers"
	"torch/torch-server/controllers/items"
	"torch/torch-server/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	testutil.TestMain(m)
}

func RouterSetup() (*httptest.ResponseRecorder, *gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)
	testutil.MockAuthMiddleware(router)
	controllers.RegisterRoutes(router, false)

	return w, c, router
}

func TestGetAllItems(t *testing.T) {
	testutil.CleanAllTables()
	testutil.SeedDB()

	// Router setup
	w, c, router := RouterSetup()

	// Getting all items
	c.Request = httptest.NewRequest(http.MethodGet, "/api/items", nil)
	router.ServeHTTP(w, c.Request)

	var items []items.Item
	if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 21, len(items))

	testutil.CleanAllTables()
}

func TestAddItem(t *testing.T) {
	// Router setup
	w, c, router := RouterSetup()

	// Request data setup
	newDreamJson := []byte(`
		{
			"title": "Test dream",
			"type": "DREAM",
			"priority": "MEDIUM",
			"targetDate": "2024-02-01"
		}
	`)
	newGoalJson := []byte(`
		{
			"title": "Test goal",
			"type": "GOAL",
			"targetDate": "2024-01-01",
			"priority": "HIGH",
			"parentID": 1
		}
	`)
	newTaskJson := []byte(`
		{
			"title": "Test task",
			"type": "TASK",
			"targetDate": "2024-01-01",
			"priority": "HIGH",
			"duration": 36000,
			"parentID": 2
		}
	`)

	var newDream items.Item
	err := json.Unmarshal(newDreamJson, &newDream)
	if err != nil {
		t.Fatal(err)
	}

	var newGoal items.Item
	err = json.Unmarshal(newGoalJson, &newGoal)
	if err != nil {
		t.Fatal(err)
	}

	var newTask items.Item
	err = json.Unmarshal(newTaskJson, &newTask)
	if err != nil {
		t.Fatal(err)
	}

	// Adding a new dream
	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-item/dream", bytes.NewReader(newDreamJson))
	router.ServeHTTP(w, c.Request)

	var returnedDream items.Item
	if err := json.Unmarshal(w.Body.Bytes(), &returnedDream); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, newDream.Title, returnedDream.Title)
	assert.Equal(t, newDream.Type, returnedDream.Type)
	assert.Equal(t, newDream.TargetDate, returnedDream.TargetDate)
	assert.Equal(t, newDream.Priority, returnedDream.Priority)

	// Router setup
	w, c, router = RouterSetup()

	// Adding a new goal
	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-item/goal", bytes.NewReader(newGoalJson))
	router.ServeHTTP(w, c.Request)

	var returnedGoal items.Item
	if err := json.Unmarshal(w.Body.Bytes(), &returnedGoal); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, returnedGoal.Title, returnedGoal.Title)
	assert.Equal(t, returnedGoal.Type, returnedGoal.Type)
	assert.Equal(t, returnedGoal.TargetDate, returnedGoal.TargetDate)
	assert.Equal(t, returnedGoal.Priority, returnedGoal.Priority)
	assert.Equal(t, returnedGoal.ParentID, returnedGoal.ParentID)

	// Router setup
	w, c, router = RouterSetup()

	// Adding a new task
	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-item/task", bytes.NewReader(newTaskJson))
	router.ServeHTTP(w, c.Request)

	var returnedTask items.Item
	if err := json.Unmarshal(w.Body.Bytes(), &returnedTask); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, newTask.Title, returnedTask.Title)
	assert.Equal(t, newTask.Type, returnedTask.Type)
	assert.Equal(t, newTask.TargetDate, returnedTask.TargetDate)
	assert.Equal(t, newTask.Priority, returnedTask.Priority)
	assert.Equal(t, newTask.ParentID, returnedTask.ParentID)
	assert.Equal(t, newTask.Duration, returnedTask.Duration)
	assert.Equal(t, newTask.Recurring, returnedTask.Recurring)

	testutil.CleanAllTables()
}

func TestUpdateItem(t *testing.T) {
	// Router setup
	w, c, router := RouterSetup()

	// Request data setup
	dreamJson := []byte(`
		{
			"title": "Test dream",
			"type": "DREAM",
			"priority": "MEDIUM",
			"targetDate": "2024-02-01"
		}
	`)

	var dream items.Item
	err := json.Unmarshal(dreamJson, &dream)
	if err != nil {
		t.Fatal(err)
	}

	// Adding a new dream
	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-item/dream", bytes.NewReader(dreamJson))
	router.ServeHTTP(w, c.Request)

	var returnedDream items.Item
	if err := json.Unmarshal(w.Body.Bytes(), &returnedDream); err != nil {
		t.Fatal(err)
	}

	// Updating the added dream
	updatedJson := []byte(`
		{
			"itemID": 1,
			"title": "New title",
			"priority": "HIGH",
			"targetDate": "2024-02-01"
		}
	`)

	var updatedDream items.ExistingDream
	err = json.Unmarshal(updatedJson, &updatedDream)
	if err != nil {
		t.Fatal(err)
	}

	// Router setup
	w, c, router = RouterSetup()

	c.Request = httptest.NewRequest(http.MethodPut, "/api/update-item/dream", bytes.NewReader(updatedJson))
	router.ServeHTTP(w, c.Request)

	if err := json.Unmarshal(w.Body.Bytes(), &returnedDream); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, updatedDream.Title, returnedDream.Title)
	assert.Equal(t, updatedDream.TargetDate, returnedDream.TargetDate)
	assert.Equal(t, updatedDream.Priority, returnedDream.Priority)

	testutil.CleanAllTables()
}

func TestRemoveItem(t *testing.T) {
	testutil.CleanAllTables()
	testutil.SeedDB()

	// Router setup
	w, c, router := RouterSetup()

	// Getting all items
	c.Request = httptest.NewRequest(http.MethodGet, "/api/items", nil)
	router.ServeHTTP(w, c.Request)

	var readItems []items.Item
	if err := json.Unmarshal(w.Body.Bytes(), &readItems); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 21, len(readItems))

	// Router setup
	w, c, router = RouterSetup()

	// Reading an item
	itemID := 1
	c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/item/%d", itemID), nil)
	router.ServeHTTP(w, c.Request)

	var item items.Item
	if err := json.Unmarshal(w.Body.Bytes(), &item); err != nil {
		t.Fatal(err)
	}

	// Router setup
	w, c, router = RouterSetup()

	// Removing the read item
	c.Request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/remove-item/%d", itemID), nil)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)

	// Router setup
	w, c, router = RouterSetup()

	// Reading the same item again
	c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/item/%d", itemID), nil)
	router.ServeHTTP(w, c.Request)

	var test items.Item
	if err := json.Unmarshal(w.Body.Bytes(), &test); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)

	testutil.CleanAllTables()
}
