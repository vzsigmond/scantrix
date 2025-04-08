package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestWatcher(t *testing.T) {
	// Create a temporary directory for testing
	dir := t.TempDir()

	// Initialize fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	// Add the temporary directory to the watcher
	err = watcher.Add(dir)
	if err != nil {
		t.Fatalf("Error adding watcher: %v", err)
	}

	// Channel to capture events
	eventChan := make(chan fsnotify.Event, 10)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				eventChan <- event
			case err := <-watcher.Errors:
				t.Logf("Watcher error: %v", err)
			}
		}
	}()

	// Set up a test file in the temporary directory
	testFile := filepath.Join(dir, "test.txt")

	// Step 1: Create the file
	err = os.WriteFile(testFile, []byte("Hello, World!"), 0644)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	select {
	case event := <-eventChan:
		t.Logf("Captured event: %v", event.Op)
		if event.Op&fsnotify.Create == 0 {
			t.Errorf("Expected CREATE event, but got: %v", event.Op)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timed out waiting for CREATE event")
	}

	// Step 2: Modify the file
	err = os.WriteFile(testFile, []byte("Updated content"), 0644)
	if err != nil {
		t.Fatalf("Error modifying test file: %v", err)
	}
	select {
	case event := <-eventChan:
		t.Logf("Captured event: %v", event.Op)
		if event.Op&fsnotify.Write == 0 {
			t.Errorf("Expected WRITE event, but got: %v", event.Op)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timed out waiting for WRITE event")
	}

	// Step 3: Delete the file
	// Allow some buffer time to handle delayed REMOVE events
	time.Sleep(100 * time.Millisecond)
	err = os.Remove(testFile)
	if err != nil {
		t.Fatalf("Error deleting test file: %v", err)
	}

	select {
	case event := <-eventChan:
		t.Logf("Captured event: %v", event.Op)
		if event.Op&fsnotify.Remove != 0 {
			t.Log("Detected REMOVE event")
		} else if event.Op&fsnotify.Write != 0 {
			t.Log("Detected WRITE event before REMOVE")
		} else {
			t.Errorf("Expected REMOVE event, but got: %v", event.Op)
		}
	case <-time.After(2 * time.Second): // Extended timeout for REMOVE
		t.Error("Timed out waiting for REMOVE event")
	}
}