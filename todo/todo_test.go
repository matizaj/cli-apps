package todo

import (
	"testing"
)

func TestAddingTaskToList(t *testing.T) {
	t.Parallel()

	// arrange
	var list List
	
	// act
	list.Add("task 1")

	// assert
	if len(list) != 1 {
		t.Errorf("got %d but expected 1", len(list))
	}
}