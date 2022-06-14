package models

import (
	"testing"
	"time"
)

var tags []string = []string{"tag1", "tag2"}

func TestTasksSerializers(t *testing.T) {
	// тест сериализации
	var tasks Tasks
	tasks = append(tasks, Task{
		Id:   1,
		Text: "first",
		Tags: tags,
		Date: JsonDate(time.Date(2009, time.November, 10, 0, 0, 0, 0, time.UTC)),
		Done: false,
	})

	js, err := tasks.Serialize()
	if err != nil {
		t.Errorf("Error marshalling: %v", err)
	}
	var tasksCheck Tasks
	tasksCheck, err = tasksCheck.Deserialize(js)
	if err != nil {
		t.Errorf("Error serialize: %v", err)
	}
	if tasksCheck[0].Id != 1 {
		t.Errorf("Error id not equal")
	}

	if tasksCheck[0].Date.String() != tasks[0].Date.String() {
		t.Errorf("Date not valid: %v", tasksCheck[0].Date.String())
	}
	if tasksCheck[0].Tags[0] != tags[0] {
		t.Errorf("Tags not valid: %v", tasksCheck[0].Tags[0])
	}
}
