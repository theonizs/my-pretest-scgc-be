package calculator

import (
	"my-go-app/internal/model"
	"testing"
)

func TestCalculateOverheat(t *testing.T) {
	// กำหนด Test Cases
	tests := []struct {
		name     string           // ชื่อเคส
		input    []model.HeatData // ข้อมูลนำเข้า
		wantNil  bool             // คาดหวังว่าเป็น null (nil) หรือไม่
		expected float64          // ค่าที่คาดหวัง (ถ้าไม่ nil)
	}{
		{
			name:    "ข้อมูลไม่พอ (น้อยกว่า 2 จุด)",
			input:   []model.HeatData{{Temp: 30}},
			wantNil: true,
		},
		{
			name: "อุณหภูมิลดลง (Slope ติดลบ)",
			input: []model.HeatData{
				{Temp: 40}, {Temp: 39}, {Temp: 38},
			},
			wantNil: true,
		},
		{
			name: "อุณหภูมิคงที่ (Slope เป็น 0)",
			input: []model.HeatData{
				{Temp: 40}, {Temp: 40}, {Temp: 40},
			},
			wantNil: true,
		},
		{
			name: "อุณหภูมิเพิ่มขึ้น (Predict ได้)",
			// สมมติ: 0->30, 1->35, 2->40 (เพิ่มทีละ 5)
			// y = 5x + 30
			// 100 = 5x + 30 => 70 = 5x => x = 14
			// ข้อมูลมี 3 จุด (Index 2)
			// อีกกี่ชั่วโมง = 14 - 3 = 11 ชั่วโมง
			input: []model.HeatData{
				{Temp: 30}, {Temp: 35}, {Temp: 40},
			},
			wantNil:  false,
			expected: 12.0, 
		},
		{
			name: "อุณหภูมิเกิน 100 ไปแล้ว (Overheated Already)",
			// 0->95, 1->105 (เพิ่มทีละ 10)
			// y = 10x + 85 (x เริ่ม 1)
			// 100 = 10x + 85 => 15 = 10x => x = 1.5
			// ปัจจุบัน N=2
			// Remaining = 1.5 - 2 = -0.5 (ติดลบ) -> ต้อง return 0
			input: []model.HeatData{
				{Temp: 95}, {Temp: 105},
			},
			wantNil:  false,
			expected: 0.0,
		},
		{
			name: "ข้อมูลเกิน 10 จุด (ตัดเอาเฉพาะ 10 จุดล่าสุด)",
			// สร้างข้อมูล 15 จุด (เกิน 10) โดยให้เป็นเส้นตรง y = 5x
			// จุดที่ 1-5: 5, 10, 15, 20, 25 (จะเป็นส่วนที่โดนตัดทิ้ง)
			// จุดที่ 6-15 (10 จุดหลัง): 30, 35, 40, ..., 75 (ส่วนที่ถูกนำมาคำนวณ)
			// การคำนวณ:
			// จุดสุดท้ายคือ 75 (ที่เวลา N=10 ของ array ใหม่)
			// เป้าหมาย 100. ขาดอีก 25.
			// Slope = 5. ดังนั้นต้องใช้เวลาอีก 25/5 = 5 ชั่วโมง
			input: func() []model.HeatData {
				data := []model.HeatData{}
				for i := 1; i <= 15; i++ {
					data = append(data, model.HeatData{Temp: float64(i * 5)})
				}
				return data
			}(),
			wantNil:  false,
			expected: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateOverheat(tt.input)

			// เช็คกรณีคาดหวัง Nil
			if tt.wantNil {
				if got != nil {
					t.Errorf("ต้องการ nil แต่ได้ %v", *got)
				}
				return
			}

			// เช็คกรณีคาดหวังค่าตัวเลข
			if got == nil {
				t.Errorf("ต้องการค่าตัวเลข แต่ได้ nil")
				return
			}

			// เช็คความแม่นยำ (เนื่องจาก float อาจคลาดเคลื่อนได้นิดหน่อย เราจึงเช็คระยะห่างแทน)
			diff := *got - tt.expected
			if diff < -0.001 || diff > 0.001 {
				t.Errorf("คำนวณผิด ได้ %v แต่ต้องการ %v", *got, tt.expected)
			}
		})
	}
}