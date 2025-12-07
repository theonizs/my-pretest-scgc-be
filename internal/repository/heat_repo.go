package repository

import (
	"my-go-app/internal/model"
	"sync"
	"time"
)

type HeatRepository interface {
	FindAll() []model.HeatData
	Save(data model.HeatData) model.HeatData
}

// memoryHeatRepository คือตัวที่ implement interface ด้วย in-memory slice
type memoryHeatRepository struct {
	db     []model.HeatData
	nextID int
	mu     sync.Mutex
}

func NewHeatRepository() HeatRepository {
	// Seed Data เบื้องต้น
	return &memoryHeatRepository{
		db: []model.HeatData{
			{ID: 1, Temp: 30.5, Time: "08", CreatedAt: time.Now().Format(time.RFC3339)},
			{ID: 2, Temp: 35.5, Time: "09", CreatedAt: time.Now().Format(time.RFC3339)},
		},
		nextID: 3,
	}
}

func (r *memoryHeatRepository) FindAll() []model.HeatData {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Return copy เพื่อความปลอดภัย (Optional)
	return r.db
}

func (r *memoryHeatRepository) Save(data model.HeatData) model.HeatData {
	r.mu.Lock()
	defer r.mu.Unlock()

	data.ID = r.nextID
	r.nextID++
	data.CreatedAt = time.Now().Format(time.RFC3339)

	r.db = append(r.db, data)
	return data
}