package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"torch/torch-server/controllers/items"
	"torch/torch-server/models"
	"torch/torch-server/testutil"

	"github.com/stretchr/testify/assert"
)

func TestGetEmptyList(t *testing.T) {
	testutil.CleanAllTables()

	// Router setup
	w, c, router := RouterSetup(userID)

	// Getting all items
	c.Request = httptest.NewRequest(http.MethodGet, "/api/items", nil)
	router.ServeHTTP(w, c.Request)

	var items []models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 0, len(items))

	testutil.CleanAllTables()
}

func TestGetAllItems(t *testing.T) {
	testutil.CleanAllTables()
	testutil.SeedDB()

	// Router setup
	w, c, router := RouterSetup(userID)

	// Getting all items
	c.Request = httptest.NewRequest(http.MethodGet, "/api/items", nil)
	router.ServeHTTP(w, c.Request)

	var items []models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 21, len(items))

	testutil.CleanAllTables()
}

func TestAddItem(t *testing.T) {
	testutil.CleanAllTables()

	// Router setup
	w, c, router := RouterSetup(userID)

	// Adding a new dream
	newDreamJson := []byte(`
		{
			"title": "Test dream",
			"type": "DREAM",
			"priority": "MEDIUM",
			"targetDate": "2024-02-01"
		}
	`)

	var newDream models.Item
	err := json.Unmarshal(newDreamJson, &newDream)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-item/dream", bytes.NewReader(newDreamJson))
	router.ServeHTTP(w, c.Request)

	var returnedDream models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &returnedDream); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, newDream.Title, returnedDream.Title)
	assert.Equal(t, newDream.Type, returnedDream.Type)
	assert.Equal(t, newDream.TargetDate, returnedDream.TargetDate)
	assert.Equal(t, newDream.Priority, returnedDream.Priority)

	// Router setup
	w, c, router = RouterSetup(userID)

	// Adding a new goal
	newGoalJson := []byte(fmt.Sprintf(`
		{
			"title": "Test goal",
			"type": "GOAL",
			"targetDate": "2024-01-01",
			"priority": "HIGH",
			"parentID": "%s"
		}
	`, returnedDream.PublicItemID))

	var newGoal models.Item
	err = json.Unmarshal(newGoalJson, &newGoal)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-item/goal", bytes.NewReader(newGoalJson))
	router.ServeHTTP(w, c.Request)

	var returnedGoal models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &returnedGoal); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, newGoal.Title, returnedGoal.Title)
	assert.Equal(t, newGoal.Type, returnedGoal.Type)
	assert.Equal(t, newGoal.TargetDate, returnedGoal.TargetDate)
	assert.Equal(t, newGoal.Priority, returnedGoal.Priority)
	assert.Equal(t, newGoal.ParentID, returnedGoal.ParentID)

	// Router setup
	w, c, router = RouterSetup(userID)

	// Adding a new task
	newTaskJson := []byte(fmt.Sprintf(`
		{
			"title": "Test task",
			"type": "TASK",
			"targetDate": "2024-01-01",
			"priority": "HIGH",
			"duration": 36000,
			"parentID": "%s"
		}
	`, returnedGoal.PublicItemID))

	var newTask models.Item
	err = json.Unmarshal(newTaskJson, &newTask)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-item/task", bytes.NewReader(newTaskJson))
	router.ServeHTTP(w, c.Request)

	var returnedTask models.Item
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
	w, c, router := RouterSetup(userID)

	// Request data setup
	dreamJson := []byte(`
		{
			"title": "Test dream",
			"type": "DREAM",
			"priority": "MEDIUM",
			"targetDate": "2024-02-01"
		}
	`)

	var dream models.Item
	err := json.Unmarshal(dreamJson, &dream)
	if err != nil {
		t.Fatal(err)
	}

	// Adding a new dream
	c.Request = httptest.NewRequest(http.MethodPost, "/api/add-item/dream", bytes.NewReader(dreamJson))
	router.ServeHTTP(w, c.Request)

	var returnedDream models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &returnedDream); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	// Router setup
	w, c, router = RouterSetup(userID)

	// Updating the added dream
	updatedJson := []byte(fmt.Sprintf(`
		{
			"itemID": "%s",
			"title": "New title",
			"priority": "HIGH",
			"targetDate": "2024-02-01"
		}
	`, returnedDream.PublicItemID))

	var updatedDream models.ExistingDream
	err = json.Unmarshal(updatedJson, &updatedDream)
	if err != nil {
		t.Fatal(err)
	}

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
	w, c, router := RouterSetup(userID)

	// Reading an item
	var publicItemID string = "1ax1usfu2uku"
	c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/item/%s", publicItemID), nil)
	router.ServeHTTP(w, c.Request)

	var item models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &item); err != nil {
		t.Fatal(err)
	}

	// Router setup
	w, c, router = RouterSetup(userID)

	// Counting the children of the item
	c.Request = httptest.NewRequest(http.MethodGet, "/api/items", nil)
	router.ServeHTTP(w, c.Request)

	var readItems []models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &readItems); err != nil {
		t.Fatal(err)
	}

	beforeChildrenNum := 0
	for _, item := range readItems {
		if item.ParentID.IsValid && item.ParentID.Val == publicItemID {
			beforeChildrenNum += 1
		}
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 21, len(readItems))
	assert.NotEqual(t, 0, beforeChildrenNum)

	// Router setup
	w, c, router = RouterSetup(userID)

	// Removing the read item
	c.Request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/remove-item/%s", publicItemID), nil)
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)

	// Router setup
	w, c, router = RouterSetup(userID)

	// Reading the same item again
	c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/item/%s", publicItemID), nil)
	router.ServeHTTP(w, c.Request)

	var test models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &test); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Router setup
	w, c, router = RouterSetup(userID)

	// Counting the children of the item
	c.Request = httptest.NewRequest(http.MethodGet, "/api/items", nil)
	router.ServeHTTP(w, c.Request)

	if err := json.Unmarshal(w.Body.Bytes(), &readItems); err != nil {
		t.Fatal(err)
	}

	afterChildrenNum := 0
	for _, item := range readItems {
		if item.ParentID.IsValid && item.ParentID.Val == publicItemID {
			afterChildrenNum += 1
		}
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 20, len(readItems))
	assert.Equal(t, 0, afterChildrenNum)

	testutil.CleanAllTables()
}

func TestUpdateItemProgress(t *testing.T) {
	testutil.CleanAllTables()
	testutil.SeedDB()

	// Router setup
	userID = uint64(1)
	w, c, router := RouterSetup(userID)

	// Reading an item
	var publicItemID string = "1ax1usfu2uku"
	c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/item/%s", publicItemID), nil)
	router.ServeHTTP(w, c.Request)

	var item models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &item); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)

	// Router setup
	w, c, router = RouterSetup(userID)

	// Updating item progress
	requestJson := []byte(fmt.Sprintf(`
		{
			"itemID": "%s",
			"timeSpent": 1000
		}
	`, item.PublicItemID))

	var requestBody items.UpdateItemProgressReq
	err := json.Unmarshal(requestJson, &requestBody)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = httptest.NewRequest(http.MethodPut, "/api/update-item-progress", bytes.NewReader(requestJson))
	router.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)

	// Router setup
	w, c, router = RouterSetup(userID)

	// Reading an item
	c.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/item/%s", publicItemID), nil)
	router.ServeHTTP(w, c.Request)

	var updatedItem models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &updatedItem); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, item.TimeSpent+requestBody.TimeSpent, updatedItem.TimeSpent)

	testutil.CleanAllTables()
}
