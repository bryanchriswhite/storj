// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package pool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFoo struct {
	called string
}

func TestGet(t *testing.T) {
	cases := []struct {
		pool          ConnectionPool
		key           string
		expected      TestFoo
		expectedError error
	}{
		{
			pool:          ConnectionPool{cache: map[string]interface{}{"foo": TestFoo{called: "hoot"}}},
			key:           "foo",
			expected:      TestFoo{called: "hoot"},
			expectedError: nil,
		},
	}

	for _, v := range cases {
		test, err := v.pool.Get(context.Background(), v.key)
		assert.Equal(t, v.expectedError, err)
		assert.Equal(t, v.expected, test)
	}
}

func TestAdd(t *testing.T) {
	cases := []struct {
		pool          ConnectionPool
		key           string
		value         TestFoo
		expected      TestFoo
		expectedError error
	}{
		{
			pool:          ConnectionPool{cache: map[string]interface{}{}},
			key:           "foo",
			value:         TestFoo{called: "hoot"},
			expected:      TestFoo{called: "hoot"},
			expectedError: nil,
		},
	}

	for _, v := range cases {
		err := v.pool.Add(context.Background(), v.key, v.value)
		assert.Equal(t, v.expectedError, err)

		test, err := v.pool.Get(context.Background(), v.key)
		assert.Equal(t, v.expectedError, err)

		assert.Equal(t, v.expected, test)
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		pool          ConnectionPool
		key           string
		value         TestFoo
		expected      interface{}
		expectedError error
	}{
		{
			pool:          ConnectionPool{cache: map[string]interface{}{"foo": TestFoo{called: "hoot"}}},
			key:           "foo",
			expected:      nil,
			expectedError: nil,
		},
	}

	for _, v := range cases {
		err := v.pool.Remove(context.Background(), v.key)
		assert.Equal(t, v.expectedError, err)

		test, err := v.pool.Get(context.Background(), v.key)
		assert.Equal(t, v.expectedError, err)

		assert.Equal(t, v.expected, test)
	}
}
