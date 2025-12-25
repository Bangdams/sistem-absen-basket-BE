package model

type DashboardResponse struct {
	TotalStudent         int     `json:"total_student"`
	PresentToday         int     `json:"present_today"`
	Absent               int     `json:"absent"`
	AttendancePercentage float64 `json:"attendance_percentage"`
}

type Result struct {
	Status string
	Total  int64
}
