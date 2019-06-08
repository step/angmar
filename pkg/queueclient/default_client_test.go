package queueclient_test

import (
	"testing"

	q "github.com/step/angmar/pkg/queueclient"
)

func TestInMemoryClient(t *testing.T) {
	client := q.NewDefaultClient()
	if err := client.Enqueue("queue", "first"); err != nil {
		t.Errorf("Unexpected error enqueueing %s in %s\n", "first", "queue")
	}

	val, err := client.Dequeue("queue")

	if err != nil {
		t.Errorf("Unexpected error dequeuing %s", "queue")
	}

	if val != "first" {
		t.Errorf("Expected %s. Got %s while dequeuing %s", "first", val, "queue")
	}

	val, err = client.Dequeue("queue")

	if (err == nil) || (val != "") {
		t.Errorf("Expected %s to be empty and expected an error", "queue")
	}
}

func TestSwitchQueueForInMemoryClient(t *testing.T) {
	client := q.NewDefaultClient()
	if err := client.Enqueue("queue", "first"); err != nil {
		t.Errorf("Unexpected error enqueueing %s in %s\n", "first", "queue")
	}

	val, err := client.SwitchQueue("queue", "another_queue")

	if err != nil {
		t.Errorf("Unexpected error dequeuing %s", "queue")
	}

	if val != "first" {
		t.Errorf("Expected %s. Got %s while dequeuing %s", "first", val, "queue")
	}

	val, err = client.Dequeue("another_queue")

	if err != nil {
		t.Errorf("Unexpected error dequeuing %s", "another_queue")
	}

	if val != "first" {
		t.Errorf("Expected %s. Got %s while dequeuing %s", "first", val, "another_queue")
	}

	val, err = client.SwitchQueue("queue", "another_queue")

	if (err == nil) || (val != "") {
		t.Errorf("Expected %s to be empty and expected an error", "queue")
	}
}
