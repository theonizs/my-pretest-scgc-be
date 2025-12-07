package calculator

import (
	"my-go-app/internal/model"
)

// CalculateOverheat ทำหน้าที่เหมือน overheatCalculator ใน TypeScript เป๊ะๆ
// Return *float64 เพื่อให้ return nil ได้ (กรณีคำนวณไม่ได้)
func CalculateOverheat(data []model.HeatData) *float64 {
	const (
		TargetTemp = 100.0
		NPoints    = 10
	)

	// 1. เลือกใช้ข้อมูล 10 จุดล่าสุด
	n := len(data)
	if n > NPoints {
		data = data[n-NPoints:]
		n = NPoints
	}

	// ตรวจสอบว่ามีข้อมูลเพียงพอหรือไม่ (อย่างน้อย 2 จุด)
	if n < 2 {
		return nil
	}

	// 2. คำนวณค่ารวม (Sum)
	var (
		sumX  float64
		sumY  float64
		sumXY float64
		sumX2 float64
	)

	// Loop คำนวณ
	// เราใช้ index+1 เป็นค่า x เหมือนใน TS (1, 2, 3, ... N)
	for i, point := range data {
		x := float64(i + 1)
		y := point.Temp

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	nf := float64(n) // แปลง n เป็น float เพื่อใช้คำนวณ

	// 3. คำนวณ Slope (m)
	denominatorM := (nf * sumX2) - (sumX * sumX)

	m := ((nf * sumXY) - (sumX * sumY)) / denominatorM

	// 4. คำนวณ Intercept (c)
	c := (sumY - (m * sumX)) / nf

	// 5. ตรวจสอบแนวโน้ม (ถ้า m <= 0 แปลว่าอุณหภูมิไม่เพิ่มขึ้น)
	if m <= 0 {
		return nil
	}

	// 6. ทำนายเวลาที่อุณหภูมิจะถึง 100°C
	// x = (TARGET - c) / m
	predictedX := (TargetTemp - c) / m

	// 7. คำนวณชั่วโมงที่เหลือ
	// ชั่วโมงปัจจุบันคือ nf (เช่น 10)
	hoursRemaining := predictedX - nf

	if hoursRemaining <= 0 {
		zero := 0.0
		return &zero
	}

	return &hoursRemaining
}