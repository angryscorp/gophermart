package utils

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func Test_Slices_Map(t *testing.T) {
	t.Run("IntToString", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []string{"1", "2", "3", "4", "5"}

		result := Map(input, func(i int) string {
			return strconv.Itoa(i)
		})

		assert.Equal(t, expected, result)
	})

	t.Run("StringToUpper", func(t *testing.T) {
		input := []string{"hello", "world", "test"}
		expected := []string{"HELLO", "WORLD", "TEST"}

		result := Map(input, strings.ToUpper)

		assert.Equal(t, expected, result)
	})

	t.Run("DoubleNumbers", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		expected := []int{2, 4, 6, 8}

		result := Map(input, func(i int) int {
			return i * 2
		})

		assert.Equal(t, expected, result)
	})

	t.Run("SingleElement", func(t *testing.T) {
		input := []int{42}
		expected := []string{"42"}

		result := Map(input, func(i int) string {
			return strconv.Itoa(i)
		})

		assert.Equal(t, expected, result)
		assert.Len(t, result, 1)
	})

	t.Run("BooleanMapping", func(t *testing.T) {
		input := []int{0, 1, 2, 0, 5}
		expected := []bool{false, true, true, false, true}

		result := Map(input, func(i int) bool {
			return i != 0
		})

		assert.Equal(t, expected, result)
	})

	t.Run("StructMapping", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		input := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}
		expected := []string{"Alice", "Bob"}

		result := Map(input, func(p Person) string {
			return p.Name
		})

		assert.Equal(t, expected, result)
	})
}
