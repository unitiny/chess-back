package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		arg1 *BasicStatus
		arg2 []int
		want int
	}{
		{
			name: "-target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: []int{-1},
			want: -9,
		},
		{
			name: "empty target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: []int{},
			want: 8,
		},
		{
			name: "one target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: []int{3},
			want: 11,
		},
		{
			name: "one target1",
			arg1: &BasicStatus{
				Status: 3,
			},
			arg2: []int{1},
			want: 2,
		},
		{
			name: "many targets",
			arg1: &BasicStatus{
				Status: 17,
			},
			arg2: []int{3, 5},
			want: 23,
		},
	}

	for i := 0; i < len(tests); i++ {
		t.Run(tests[i].name, func(t *testing.T) {
			tests[i].arg1.Reverse(tests[i].arg2...)
			assert.Equal(t, tests[i].arg1.Status, tests[i].want)
		})
	}
}

func TestReset(t *testing.T) {
	tests := []struct {
		name string
		arg1 *BasicStatus
		arg2 []int
		want int
	}{
		{
			name: "-target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: []int{-1},
			want: 0,
		},
		{
			name: "empty target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: []int{},
			want: 8,
		},
		{
			name: "one target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: []int{3},
			want: 8,
		},
		{
			name: "one target1",
			arg1: &BasicStatus{
				Status: 3,
			},
			arg2: []int{1},
			want: 2,
		},
		{
			name: "many targets",
			arg1: &BasicStatus{
				Status: 17,
			},
			arg2: []int{3, 5},
			want: 16,
		},
	}

	for i := 0; i < len(tests); i++ {
		t.Run(tests[i].name, func(t *testing.T) {
			tests[i].arg1.Reset(tests[i].arg2...)
			assert.Equal(t, tests[i].want, tests[i].arg1.Status)
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name string
		arg1 *BasicStatus
		arg2 []int
		want int
	}{
		{
			name: "-target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: []int{-2},
			want: -2,
		},
		{
			name: "empty target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: []int{},
			want: 8,
		},
		{
			name: "one target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: []int{3},
			want: 11,
		},
		{
			name: "one target1",
			arg1: &BasicStatus{
				Status: 3,
			},
			arg2: []int{0},
			want: 3,
		},
		{
			name: "many targets",
			arg1: &BasicStatus{
				Status: 17,
			},
			arg2: []int{3, 5},
			want: 23,
		},
		{
			name: "many targets1",
			arg1: &BasicStatus{
				Status: 0,
			},
			arg2: []int{1, 2},
			want: 3,
		},
	}

	for i := 0; i < len(tests); i++ {
		t.Run(tests[i].name, func(t *testing.T) {
			tests[i].arg1.Set(tests[i].arg2...)
			assert.Equal(t, tests[i].want, tests[i].arg1.Status)
		})
	}
}

func TestHas(t *testing.T) {
	tests := []struct {
		name string
		arg1 *BasicStatus
		arg2 int
		want bool
	}{
		{
			name: "-target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: -2,
			want: true,
		},
		{
			name: "empty target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: 0,
			want: false,
		},
		{
			name: "one target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: 3,
			want: false,
		},
		{
			name: "right target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: 8,
			want: true,
		},
	}

	for i := 0; i < len(tests); i++ {
		t.Run(tests[i].name, func(t *testing.T) {
			res := tests[i].arg1.Has(tests[i].arg2)
			assert.Equal(t, tests[i].want, res)
		})
	}
}

func TestIsSame(t *testing.T) {

	tests := []struct {
		name string
		arg1 *BasicStatus
		arg2 int
		arg3 int
		want bool
	}{
		{
			name: "-target",
			arg1: &BasicStatus{
				Status: 11,
			},
			arg2: -5,
			arg3: 0,
			want: true,
		},
		{
			name: "empty target",
			arg1: &BasicStatus{
				Status: 8,
			},
			arg2: 0,
			arg3: 0,
			want: true,
		},
		{
			name: "one target",
			arg1: &BasicStatus{
				Status: 12,
			},
			arg2: 1,
			arg3: 0,
			want: false,
		},
		{
			name: "one target",
			arg1: &BasicStatus{
				Status: 11,
			},
			arg2: 3,
			arg3: 0,
			want: true,
		},
		{
			name: "offset target",
			arg1: &BasicStatus{
				Status: 11,
			},
			arg2: 2,
			arg3: 1,
			want: false,
		},
		{
			name: "offset target1",
			arg1: &BasicStatus{
				Status: 11,
			},
			arg2: 2,
			arg3: 2,
			want: true,
		},
		{
			name: "offset target2",
			arg1: &BasicStatus{
				Status: 15,
			},
			arg2: 5,
			arg3: 1,
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := test.arg1.IsSame(test.arg2, test.arg3)
			assert.Equal(t, test.want, res)
		})
	}
}
