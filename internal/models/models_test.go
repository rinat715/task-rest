package models

import (
	"testing"
	"time"
)

func TestTasksSerializers(t *testing.T) {
	// тест сериализации
	var tasks Tasks
	tasks = append(tasks, Task{
		Id:   1,
		Text: "first",
		Tags: nil,
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

	if tasksCheck[0].Date.String() == "2009-11-10" {
		t.Errorf("Date not valid: %v", tasks[0].Date.String())
	}
}
