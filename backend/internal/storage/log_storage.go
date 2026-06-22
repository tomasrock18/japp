package storage

import (
	"sync"
	"time"

	"github.com/tomasrock18/japp/backend/internal/model"
)

type LogStorage struct {
	mu     sync.RWMutex
	logs   []model.FoodLog
	nextID int64
}

func NewLogStorage() *LogStorage {
	return &LogStorage{
		logs:   make([]model.FoodLog, 0),
		nextID: 1,
	}
}

func (s *LogStorage) CreateLog(log model.FoodLog) (model.FoodLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.ID = s.nextID
	log.CreatedAt = time.Now()
	s.nextID++

	s.logs = append(s.logs, log)
	return log, nil
}

func (s *LogStorage) GetLogsByUserAndDate(telegramID int64, date time.Time) ([]model.FoodLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []model.FoodLog
	for _, log := range s.logs {
		if log.TelegramID == telegramID && isSameDay(log.CreatedAt, date) {
			result = append(result, log)
		}
	}

	if result == nil {
		result = []model.FoodLog{}
	}
	return result, nil
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

func (s *LogStorage) GetLogsByUser(telegramID int64, limit, offset int) ([]model.FoodLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []model.FoodLog
	for _, log := range s.logs {
		if log.TelegramID == telegramID {
			result = append(result, log)
		}
		if len(result) == limit {
			break
		}
	}

	if result == nil {
		result = []model.FoodLog{}
	}

	if offset > len(result) {
		return []model.FoodLog{}, nil
	}

	return result[offset:], nil
}
