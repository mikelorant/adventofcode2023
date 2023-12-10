package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPart1FurthestSteps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		want     int
	}{
		{
			name:     "demo1",
			filename: "demo1.txt",
			want:     4,
		},
		{
			name:     "demo2",
			filename: "demo2.txt",
			want:     8,
		},
		{
			name:     "input1",
			filename: "input1.txt",
			want:     6640,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			steps := FurthestSteps(tt.filename)
			assert.Equal(t, tt.want, steps)
		})
	}
}

func TestPart2EnclosedTiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		want     int
	}{
		{
			name:     "demo3",
			filename: "demo3.txt",
			want:     4,
		},
		{
			name:     "demo4",
			filename: "demo4.txt",
			want:     8,
		},
		{
			name:     "demo5",
			filename: "demo5.txt",
			want:     10,
		},
		{
			name:     "input1",
			filename: "input1.txt",
			want:     411,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tiles := EnclosedTiles(tt.filename)
			assert.Equal(t, tt.want, tiles)
		})
	}
}
