package logging

import (
	"os"
	"strings"
	"sync"
	"testing"
)

func TestLogger_WriteAndRead(t *testing.T) {
	testFile := "test_cura.log"
	defer os.Remove(testFile)

	l := NewLogger(testFile)

	// test write
	l.Write("Test message 1")
	l.Write("Test message 2")

	// test read
	content, err := l.ReadLogs(10)
	if err != nil {
		t.Fatalf("Failed to read logs: %v", err)
	}

	if !strings.Contains(content, "Test message 1") || !strings.Contains(content, "Test message 2") {
		t.Errorf("Logs do not contain expected messages. Got: %s", content)
	}
}

func TestLogger_Concurrency(t *testing.T) {
	testFile := "concurrency_test.log"
	defer os.Remove(testFile)

	l := NewLogger(testFile)
	var wg sync.WaitGroup

	// stress test: 50 for e.g goroutines writing and reading at the same time
	for i := 0; i < 50; i++ {
		wg.Add(2)
		
		// concurrent Writer
		go func(id int) {
			defer wg.Done()
			l.Write("Concurrent write from routine")
		}(i)

		// concurrent Reader
		go func() {
			defer wg.Done()
			_, _ = l.ReadLogs(10)
		}()
	}

	wg.Wait()

	content, _ := l.ReadLogs(100)
	lines := strings.Split(strings.TrimSpace(content), "\n")
	
	if len(lines) < 1 {
		t.Error("Concurrent test resulted in empty log file")
	}
}