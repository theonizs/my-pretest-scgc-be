package handler

import (
	"bytes"
	"encoding/json"
	"my-go-app/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// --- 1. ปรับปรุง Mock Repo ให้ฉลาดขึ้น ---
type mockRepo struct {
	mockData []model.HeatData
}

func (m *mockRepo) FindAll() []model.HeatData {
	return m.mockData
}

func (m *mockRepo) Save(data model.HeatData) model.HeatData {
	// จำลองการ Save แล้วคืนค่าพร้อม ID
	data.ID = 999
	return data
}

// Helper สร้างข้อมูลจำลอง
func generateMockData(count int) []model.HeatData {
	data := []model.HeatData{}
	for i := 1; i <= count; i++ {
		data = append(data, model.HeatData{ID: i, Temp: float64(i * 10), Time: "00"})
	}
	return data
}

// --- Test 1: GET ข้อมูลเยอะ (ต้อง Slice) ---
func TestGetHeatData_SliceLast10(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock มี 15 ตัว
	mockRepository := &mockRepo{mockData: generateMockData(15)}
	h := NewHeatHandler(mockRepository)

	r := gin.Default()
	r.GET("/test", h.GetHeatData)

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var rawMap map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &rawMap)
	dataList := rawMap["data"].(map[string]interface{})["data"].([]interface{})

	// ต้องเหลือแค่ 10
	if len(dataList) != 10 {
		t.Errorf("Expected 10 items, got %d", len(dataList))
	}
}

// --- Test 2: GET ข้อมูลน้อย (ไม่ต้อง Slice) ---
func TestGetHeatData_NoSlice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock มีแค่ 5 ตัว
	mockRepository := &mockRepo{mockData: generateMockData(5)}
	h := NewHeatHandler(mockRepository)

	r := gin.Default()
	r.GET("/test", h.GetHeatData)

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var rawMap map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &rawMap)
	dataList := rawMap["data"].(map[string]interface{})["data"].([]interface{})

	// ต้องได้ครบ 5 ตัว ไม่โดนตัด
	if len(dataList) != 5 {
		t.Errorf("Expected 5 items, got %d", len(dataList))
	}
}

// --- Test 3: POST สำเร็จ ---
func TestCreateHeatData_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHeatHandler(&mockRepo{}) // Mock ว่างๆ ก็ได้เพราะ Save ไม่ได้เช็คอะไรใน mock

	r := gin.Default()
	r.POST("/test", h.CreateHeatData)

	// Prepare JSON Body
	payload := model.HeatData{Temp: 50.5, Time: "12"}
	jsonValue, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d", w.Code)
	}
}

// --- Test 4: POST ส่งข้อมูลผิด (Bad Request) ---
func TestCreateHeatData_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHeatHandler(&mockRepo{})

	r := gin.Default()
	r.POST("/test", h.CreateHeatData)

	// ส่ง JSON มั่วๆ ที่แปลงเป็น float ไม่ได้ หรือ field หาย
	jsonValue := []byte(`{"temp": "not-a-number", "time": 123}`)

	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", w.Code)
	}
}