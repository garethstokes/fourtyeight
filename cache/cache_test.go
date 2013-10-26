package cache

import (
  "testing"
)

func TestInsertAndGet(t * testing.T) {
  test_value := "test_value"

  Set("test", "test_key", test_value)
  result := Get("test", "test_key")

  if result != test_value {
    t.Fail()
  }
}
