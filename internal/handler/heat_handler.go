package handler

import (
	"my-go-app/internal/logic/calculator"
	"my-go-app/internal/model"
	"my-go-app/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HeatHandler struct {
	repo repository.HeatRepository
}

func NewHeatHandler(repo repository.HeatRepository) *HeatHandler {
	return &HeatHandler{repo: repo}
}

// GET /api/heat-data
func (h *HeatHandler) GetHeatData(c *gin.Context) {
	allData := h.repo.FindAll()

	var recentData []model.HeatData
	limit := 10

	if len(allData) > limit {
		// Slice เอาตั้งแต่ตัวที่ (ความยาว - 10) ไปจนจบ
		recentData = allData[len(allData)-limit:]
	} else {
		recentData = allData
	}

	prediction := calculator.CalculateOverheat(recentData)

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data: model.HeatDataResponse{
			Data:       recentData,
			Prediction: prediction,
		},
	})
}

// POST /api/heat-data
func (h *HeatHandler) CreateHeatData(c *gin.Context) {
	var input model.HeatData

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	savedData := h.repo.Save(input)

	c.JSON(http.StatusCreated, model.Response{
		Status:  http.StatusCreated,
		Message: "Data saved successfully",
		Data:    savedData,
	})
}