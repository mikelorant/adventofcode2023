package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSumPathBetweenPairs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		filename       string
		emptySpaceSize int
		want           int
	}{
		{
			name:           "demo1",
			filename:       "demo1.txt",
			emptySpaceSize: 2,
			want:           374,
		},
		{
			name:           "input1",
			filename:       "input1.txt",
			emptySpaceSize: 2,
			want:           9769724,
		},
		{
			name:           "demo1",
			filename:       "demo1.txt",
			emptySpaceSize: 10,
			want:           1030,
		},
		{
			name:           "demo1",
			filename:       "demo1.txt",
			emptySpaceSize: 100,
			want:           8410,
		},
		{
			name:           "input1",
			filename:       "input1.txt",
			emptySpaceSize: 1000000,
			want:           603020563700,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			sum := SumPathBetweenPairs(tt.filename, tt.emptySpaceSize)
			assert.Equal(t, tt.want, sum)
		})
	}
}
