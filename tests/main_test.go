package tests

import (
	"net/http/httptest"
	"testing"
	"torch/torch-server/controllers"
	"torch/torch-server/testutil"

	"github.com/gin-gonic/gin"
)

var userID = uint64(1)

func TestMain(m *testing.M) {
	testutil.TestMain(m)
}

func RouterSetup(mockUserID uint64) (*httptest.ResponseRecorder, *gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)
	testutil.MockAuthMiddleware(router, mockUserID)
	controllers.RegisterRoutes(router, false)

	return w, c, router
}
