package model

import "sunhost/config"

type UserLog struct {
	ID         int    `json:"id"`
	Username   string `json:"username" gorm:"foreignKey:Username;references:Username"`
	IPAddress  string `json:"ip_address"`
	SystemInfo string `json:"system_info"`
	Action     string `json:"action"`
	Time       string `json:"time"`
}

func (ul *UserLog) TableName() string {
	return "logs"
}

func (ul *UserLog) Create() error {
	_, err := config.DB.Exec(`
		INSERT INTO logs (username, ip_address, system_info, action)
		VALUES (?, ?, ?, ?)
	`, ul.Username, ul.IPAddress, ul.SystemInfo, ul.Action)
	return err
}

func (ul *UserLog) GetByUsername(username string) ([]UserLog, error) {
	rows, err := config.DB.Query(`
		SELECT id, username, ip_address, system_info, action , time
		FROM logs
		WHERE username = ?
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []UserLog
	for rows.Next() {
		var log UserLog
		if err := rows.Scan(&log.ID, &log.Username, &log.IPAddress, &log.SystemInfo, &log.Action, &log.Time); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (ul *UserLog) CheckCountByUsername(username string) (int, error) {
	var count int
	err := config.DB.QueryRow(`
		SELECT COUNT(*) FROM logs WHERE username = ?`, username).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ul *UserLog) DeleteByID() error {
	_, err := config.DB.Exec(`
		DELETE FROM logs WHERE id = ? 
	`, ul.ID)
	if err != nil {
		return err
	}

	return nil
}
func (ul *UserLog) DeleteOnceByUsername() error {
    _, err := config.DB.Exec(`
        DELETE FROM logs 
        WHERE id = (
            SELECT id 
            FROM logs 
            WHERE username = ? AND action != 'User Created' 
            ORDER BY id ASC 
            LIMIT 1
        )
    `, ul.Username)
    
    if err != nil {
        return err
    }
    return nil
}