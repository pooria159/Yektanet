package controllers

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go-ad-panel/models"
    "go-ad-panel/repositories"
    "gorm.io/gorm"
)

// MockPublisherRepository is a mock implementation of PublisherRepositoryInterface
type MockPublisherRepository struct {
    mock.Mock
}

func (m *MockPublisherRepository) FindByID(id uint) (models.Publisher, error) {
    args := m.Called(id)
    return args.Get(0).(models.Publisher), args.Error(1)
}

func (m *MockPublisherRepository) Save(p *models.Publisher) error {
    args := m.Called(p)
    return args.Error(0)
}

func (m *MockPublisherRepository) Update(p *models.Publisher) error {
    args := m.Called(p)
    return args.Error(0)
}

func (m *MockPublisherRepository) Delete(id uint) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *MockPublisherRepository) FindAll() ([]models.Publisher, error) {
    args := m.Called()
    return args.Get(0).([]models.Publisher), args.Error(1)
}

var _ repositories.PublisherRepositoryInterface = (*MockPublisherRepository)(nil)

func TestPublisherPanel(t *testing.T) {
    gin.SetMode(gin.TestMode)

    t.Run("Valid ID and Publisher Exists", func(t *testing.T) {
        mockRepo := new(MockPublisherRepository)
        ctrl := PublisherController{Repo: mockRepo}
        router := gin.Default()
        router.LoadHTMLGlob("../templates/*")
        router.GET("/publisher/:id", ctrl.PublisherPanel)
        publisher := models.Publisher{Model: gorm.Model{ID: 1}, Name: "Test Publisher"}
        mockRepo.On("FindByID", uint(1)).Return(publisher, nil)
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/publisher/1", nil)
        router.ServeHTTP(w, req)
        assert.Equal(t, http.StatusOK, w.Code)
        assert.Contains(t, w.Body.String(), "Test Publisher")
    })

    t.Run("Invalid ID", func(t *testing.T) {
        mockRepo := new(MockPublisherRepository)
        ctrl := PublisherController{Repo: mockRepo}
        router := gin.Default()
        router.LoadHTMLGlob("../templates/*")
        router.GET("/publisher/:id", ctrl.PublisherPanel)

        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/publisher/abc", nil)
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)
        assert.Contains(t, w.Body.String(), "Invalid ID")
    })

    t.Run("Publisher Not Found", func(t *testing.T) {
        mockRepo := new(MockPublisherRepository)
        ctrl := PublisherController{Repo: mockRepo}
        router := gin.Default()
        router.LoadHTMLGlob("../templates/*")
        router.GET("/publisher/:id", ctrl.PublisherPanel)
        mockRepo.On("FindByID", uint(2)).Return(models.Publisher{}, fmt.Errorf("Publisher not found"))

        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/publisher/2", nil)
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusNotFound, w.Code)
        assert.Contains(t, w.Body.String(), "Publisher not found")
    })
}

// ---------------------------------------------------------------TestPublisherPanel----------------------------------------------------------------

func TestPublisherWithdraw(t *testing.T) {
    gin.SetMode(gin.TestMode)

    t.Run("Invalid ID", func(t *testing.T) {
        mockRepo := new(MockPublisherRepository)
        ctrl := PublisherController{Repo: mockRepo}
        router := gin.Default()
        router.POST("/publisher/:id/withdraw", ctrl.PublisherWithdraw)

        w := httptest.NewRecorder()
        req, _ := http.NewRequest("POST", "/publisher/abc/withdraw", nil)
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)
        assert.Contains(t, w.Body.String(), "Invalid ID")
    })

    t.Run("Publisher Not Found", func(t *testing.T) {
        mockRepo := new(MockPublisherRepository)
        ctrl := PublisherController{Repo: mockRepo}
        router := gin.Default()
        router.POST("/publisher/:id/withdraw", ctrl.PublisherWithdraw)

        mockRepo.On("FindByID", uint(2)).Return(models.Publisher{}, fmt.Errorf("Publisher not found"))

        w := httptest.NewRecorder()
        req, _ := http.NewRequest("POST", "/publisher/2/withdraw", nil)
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusNotFound, w.Code)
        assert.Contains(t, w.Body.String(), "Publisher not found")
    })

    t.Run("Invalid Amount", func(t *testing.T) {
        mockRepo := new(MockPublisherRepository)
        ctrl := PublisherController{Repo: mockRepo}
        router := gin.Default()
        router.POST("/publisher/:id/withdraw", ctrl.PublisherWithdraw)

        publisher := models.Publisher{Model: gorm.Model{ID: 1}, Credit: 100}
        mockRepo.On("FindByID", uint(1)).Return(publisher, nil)

        w := httptest.NewRecorder()
        form := url.Values{}
        form.Set("amount", "invalid")
        req, _ := http.NewRequest("POST", "/publisher/1/withdraw", strings.NewReader(form.Encode()))
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)
        assert.Contains(t, w.Body.String(), "Invalid amount")
    })

    t.Run("Insufficient Balance", func(t *testing.T) {
        mockRepo := new(MockPublisherRepository)
        ctrl := PublisherController{Repo: mockRepo}
        router := gin.Default()
        router.POST("/publisher/:id/withdraw", ctrl.PublisherWithdraw)

        publisher := models.Publisher{Model: gorm.Model{ID: 1}, Credit: 50}
        mockRepo.On("FindByID", uint(1)).Return(publisher, nil)

        w := httptest.NewRecorder()
        form := url.Values{}
        form.Set("amount", "100")
        req, _ := http.NewRequest("POST", "/publisher/1/withdraw", strings.NewReader(form.Encode()))
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)
        assert.Contains(t, w.Body.String(), "Insufficient balance")
    })

    t.Run("Successful Withdrawal", func(t *testing.T) {
        mockRepo := new(MockPublisherRepository)
        ctrl := PublisherController{Repo: mockRepo}
        router := gin.Default()
        router.POST("/publisher/:id/withdraw", ctrl.PublisherWithdraw)

        publisher := models.Publisher{Model: gorm.Model{ID: 1}, Credit: 100}
        mockRepo.On("FindByID", uint(1)).Return(publisher, nil)
        mockRepo.On("Update", mock.AnythingOfType("*models.Publisher")).Return(nil)

        w := httptest.NewRecorder()
        form := url.Values{}
        form.Set("amount", "50")
        req, _ := http.NewRequest("POST", "/publisher/1/withdraw", strings.NewReader(form.Encode()))
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusSeeOther, w.Code)
    })
}