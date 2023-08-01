package model

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	mockArgs := m.Called(value, conds)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Raw(sql string, values ...interface{}) *gorm.DB {
	mockArgs := m.Called(sql, values)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Scan(dest interface{}) *gorm.DB {
	mockArgs := m.Called(dest)
	return mockArgs.Get(0).(*gorm.DB)
}
