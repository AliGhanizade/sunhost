package model

import (
	"sunhost/config"
)

type SystemStat struct {
	ID          int     `json:"id"`
	AllocRAM    float64 `json:"alloc_ram"`
	Goroutines  int     `json:"goroutines"`
	LiveObjects uint64  `json:"live_objects"`
	RecordTime  string  `json:"record_time"`
}

func (ss *SystemStat) TableName() string {
	return "system_stats"
}

// ذخیره وضعیت سیستم در دیتابیس SQLite
func (ss *SystemStat) Create() error {
	_, err := config.DB.Exec(`
		INSERT INTO system_stats (alloc_ram, goroutines, live_objects, record_time)
		VALUES (?, ?, ?, datetime('now', 'localtime'))
	`, ss.AllocRAM, ss.Goroutines, ss.LiveObjects)
	return err
}

// نگه داشتن فقط ۱۰۰۰ رکورد آخر و پاک کردن بقیه (جلوگیری از انفجار دیتابیس)
func (ss *SystemStat) CleanupOldRecords() error {
	_, err := config.DB.Exec(`
		DELETE FROM system_stats 
		WHERE id NOT IN (
			SELECT id FROM system_stats ORDER BY id DESC LIMIT 1000
		)
	`)
	return err
}

// استراکچر کمکی برای ارسال دیتای آماری به فرانت‌اند
type SystemSummary struct {
	AvgRam      float64 `json:"avg_ram_today"`
	MaxThreads  int     `json:"max_goroutines"`
	TotalRecord int     `json:"total_records"`
}

// گرفتن میانگین‌ها و رکوردهای کلی برای باکس‌های بالای صفحه
func (ss *SystemStat) GetSummary() (SystemSummary, error) {
	var sum SystemSummary
	// COALESCE کمک می‌کند اگر دیتابیس خالی بود، به جای ارور، عدد 0 برگردد
	err := config.DB.QueryRow(`
		SELECT 
			COALESCE(AVG(alloc_ram), 0), 
			COALESCE(MAX(goroutines), 0), 
			COUNT(id) 
		FROM system_stats
	`).Scan(&sum.AvgRam, &sum.MaxThreads, &sum.TotalRecord)
	
	return sum, err
}