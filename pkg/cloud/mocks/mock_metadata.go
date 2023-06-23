// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kubernetes-sigs/aws-fsx-openzfs-csi-driver/pkg/cloud (interfaces: MetadataService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetadataService is a mock of MetadataService interface.
type MockMetadataService struct {
	ctrl     *gomock.Controller
	recorder *MockMetadataServiceMockRecorder
}

// MockMetadataServiceMockRecorder is the mock recorder for MockMetadataService.
type MockMetadataServiceMockRecorder struct {
	mock *MockMetadataService
}

// NewMockMetadataService creates a new mock instance.
func NewMockMetadataService(ctrl *gomock.Controller) *MockMetadataService {
	mock := &MockMetadataService{ctrl: ctrl}
	mock.recorder = &MockMetadataServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetadataService) EXPECT() *MockMetadataServiceMockRecorder {
	return m.recorder
}

// GetAvailabilityZone mocks base method.
func (m *MockMetadataService) GetAvailabilityZone() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailabilityZone")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAvailabilityZone indicates an expected call of GetAvailabilityZone.
func (mr *MockMetadataServiceMockRecorder) GetAvailabilityZone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailabilityZone", reflect.TypeOf((*MockMetadataService)(nil).GetAvailabilityZone))
}

// GetInstanceID mocks base method.
func (m *MockMetadataService) GetInstanceID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstanceID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetInstanceID indicates an expected call of GetInstanceID.
func (mr *MockMetadataServiceMockRecorder) GetInstanceID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstanceID", reflect.TypeOf((*MockMetadataService)(nil).GetInstanceID))
}

// GetRegion mocks base method.
func (m *MockMetadataService) GetRegion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegion")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetRegion indicates an expected call of GetRegion.
func (mr *MockMetadataServiceMockRecorder) GetRegion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegion", reflect.TypeOf((*MockMetadataService)(nil).GetRegion))
}
