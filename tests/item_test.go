package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"torch/torch-server/models"
	r "torch/torch-server/router"
	"torch/torch-server/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	testutil.TestMain(m)
}

func TestGetAllItems(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(w)
	testutil.MockAuthMiddleware(router)
	r.RegisterRoutes(router, false)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/items", nil)
	router.ServeHTTP(w, ctx.Request)

	var items []models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, len(items), 21)
}

func TestAddItem(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(w)
	testutil.MockAuthMiddleware(router)
	r.RegisterRoutes(router, false)

	jsonData := []byte(`
		{
			"title": "Test item",
			"type": "TASK",
			"targetDate": "2024-01-01",
			"priority": "HIGH",
			"duration": 36000
		}	
	`)

	var newItem models.Item
	err := json.Unmarshal(jsonData, &newItem)
	if err != nil {
		t.Fatal(err)
	}

	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/add-item", bytes.NewReader(jsonData))
	router.ServeHTTP(w, ctx.Request)

	var returnedItem models.Item
	if err := json.Unmarshal(w.Body.Bytes(), &returnedItem); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, newItem.Title, returnedItem.Title)
	assert.Equal(t, newItem.Type, returnedItem.Type)
	// assert.Equal(t, newItem.TargetDate.String(), returnedItem.TargetDate.String())
	assert.Equal(t, newItem.Priority, returnedItem.Priority)
	assert.Equal(t, newItem.Duration, returnedItem.Duration)
	// assert.Equal(t, newItem.Recurring, item.Recurring)
	assert.Equal(t, newItem.ParentID.String(), returnedItem.ParentID.String())
}
