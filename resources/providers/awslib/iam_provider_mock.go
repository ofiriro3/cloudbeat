// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
// Code generated by mockery v2.10.4. DO NOT EDIT.

package awslib

import (
	context "context"

	iam "github.com/aws/aws-sdk-go-v2/service/iam"
	mock "github.com/stretchr/testify/mock"
)

// MockIAMRolePermissionGetter is an autogenerated mock type for the IAMRolePermissionGetter type
type MockIAMRolePermissionGetter struct {
	mock.Mock
}

type MockIAMRolePermissionGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIAMRolePermissionGetter) EXPECT() *MockIAMRolePermissionGetter_Expecter {
	return &MockIAMRolePermissionGetter_Expecter{mock: &_m.Mock}
}

// GetIAMRolePermissions provides a mock function with given fields: ctx, roleName
func (_m *MockIAMRolePermissionGetter) GetIAMRolePermissions(ctx context.Context, roleName string) ([]iam.GetRolePolicyResponse, error) {
	ret := _m.Called(ctx, roleName)

	var r0 []iam.GetRolePolicyResponse
	if rf, ok := ret.Get(0).(func(context.Context, string) []iam.GetRolePolicyResponse); ok {
		r0 = rf(ctx, roleName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]iam.GetRolePolicyResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, roleName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIAMRolePermissionGetter_GetIAMRolePermissions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIAMRolePermissions'
type MockIAMRolePermissionGetter_GetIAMRolePermissions_Call struct {
	*mock.Call
}

// GetIAMRolePermissions is a helper method to define mock.On call
//  - ctx context.Context
//  - roleName string
func (_e *MockIAMRolePermissionGetter_Expecter) GetIAMRolePermissions(ctx interface{}, roleName interface{}) *MockIAMRolePermissionGetter_GetIAMRolePermissions_Call {
	return &MockIAMRolePermissionGetter_GetIAMRolePermissions_Call{Call: _e.mock.On("GetIAMRolePermissions", ctx, roleName)}
}

func (_c *MockIAMRolePermissionGetter_GetIAMRolePermissions_Call) Run(run func(ctx context.Context, roleName string)) *MockIAMRolePermissionGetter_GetIAMRolePermissions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockIAMRolePermissionGetter_GetIAMRolePermissions_Call) Return(_a0 []iam.GetRolePolicyResponse, _a1 error) *MockIAMRolePermissionGetter_GetIAMRolePermissions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}
