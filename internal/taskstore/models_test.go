package taskstore

import (
	"testing"
	"time"
)

func TestTasksSerializers(t *testing.T) {
	// тест сериализации
	ts := New()
	var d JsonDate = JsonDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
	id := ts.Create("Hola", nil, JsonDate(d), false).Id
	allTasks := ts.GetAll()
	js, err := allTasks.Serialize()
	if err != nil {
		t.Errorf("Error marshalling: %v", err)
	}
	var task Tasks
	tasks, err := task.Deserialize(js)
	if err != nil {
		t.Errorf("Error serialize: %v", err)
	}
	if tasks[0].Id != id {
		t.Errorf("Error id not equal")
	}

	if tasks[0].Date.String() == "2009-11-10" {
		t.Errorf("Date not valid: %v", tasks[0].Date.String())
	}
}
