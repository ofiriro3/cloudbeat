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

	eks "github.com/aws/aws-sdk-go-v2/service/eks"
	mock "github.com/stretchr/testify/mock"
)

// MockedEksClusterDescriber is an autogenerated mock type for the EksClusterDescriber type
type MockedEksClusterDescriber struct {
	mock.Mock
}

type MockEksClusterDescriber_Expecter struct {
	mock *mock.Mock
}

func (_m *MockedEksClusterDescriber) EXPECT() *MockEksClusterDescriber_Expecter {
	return &MockEksClusterDescriber_Expecter{mock: &_m.Mock}
}

// DescribeCluster provides a mock function with given fields: ctx, clusterName
func (_m *MockedEksClusterDescriber) DescribeCluster(ctx context.Context, clusterName string) (*eks.DescribeClusterResponse, error) {
	ret := _m.Called(ctx, clusterName)

	var r0 *eks.DescribeClusterResponse
	if rf, ok := ret.Get(0).(func(context.Context, string) *eks.DescribeClusterResponse); ok {
		r0 = rf(ctx, clusterName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*eks.DescribeClusterResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, clusterName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEksClusterDescriber_DescribeCluster_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DescribeCluster'
type MockEksClusterDescriber_DescribeCluster_Call struct {
	*mock.Call
}

// DescribeCluster is a helper method to define mock.On call
//  - ctx context.Context
//  - clusterName string
func (_e *MockEksClusterDescriber_Expecter) DescribeCluster(ctx interface{}, clusterName interface{}) *MockEksClusterDescriber_DescribeCluster_Call {
	return &MockEksClusterDescriber_DescribeCluster_Call{Call: _e.mock.On("DescribeCluster", ctx, clusterName)}
}

func (_c *MockEksClusterDescriber_DescribeCluster_Call) Run(run func(ctx context.Context, clusterName string)) *MockEksClusterDescriber_DescribeCluster_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockEksClusterDescriber_DescribeCluster_Call) Return(_a0 *eks.DescribeClusterResponse, _a1 error) *MockEksClusterDescriber_DescribeCluster_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}
