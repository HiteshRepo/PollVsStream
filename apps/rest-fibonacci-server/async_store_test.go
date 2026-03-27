package restfibonacciserver

import (
	"sync"
	"testing"
)

func TestAsyncStoreWriteAndRead(t *testing.T) {
	store := NewAsyncStore(5)
	store.Write(0, 0)
	store.Write(1, 1)

	nums, current, requested := store.Read()
	if len(nums) != 2 {
		t.Fatalf("expected 2 numbers, got %d", len(nums))
	}
	if nums[0] != 0 || nums[1] != 1 {
		t.Errorf("unexpected numbers: %v", nums)
	}
	if current != 1 {
		t.Errorf("expected current=1, got %d", current)
	}
	if requested != 5 {
		t.Errorf("expected requested=5, got %d", requested)
	}
}

func TestAsyncStoreDrainOnRead(t *testing.T) {
	store := NewAsyncStore(3)
	store.Write(0, 0)

	nums, _, _ := store.Read()
	if len(nums) != 1 {
		t.Fatalf("first read: expected 1 number, got %d", len(nums))
	}

	nums2, _, _ := store.Read()
	if len(nums2) != 0 {
		t.Errorf("second read: expected 0 numbers (drained), got %d", len(nums2))
	}
}

func TestAsyncStoreConcurrentWrites(t *testing.T) {
	store := NewAsyncStore(100)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			store.Write(n, n)
		}(i)
	}
	wg.Wait()

	nums, current, _ := store.Read()
	if len(nums) != 100 {
		t.Errorf("expected 100 numbers after concurrent writes, got %d", len(nums))
	}
	if current < 0 || current >= 100 {
		t.Errorf("current out of range: %d", current)
	}
}

func TestAsyncStoreEndOfSequence(t *testing.T) {
	store := NewAsyncStore(3)
	store.Write(0, 0)
	store.Write(1, 1)
	store.Write(1, 2)

	_, current, requested := store.Read()
	if (current + 1) == requested {
		// end of sequence detected correctly
		return
	}
	t.Errorf("expected end of sequence: current=%d requested=%d", current, requested)
}
