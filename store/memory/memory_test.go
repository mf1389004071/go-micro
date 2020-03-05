package memory

import (
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/micro/go-micro/v2/store"
)

func TestMemoryReInit(t *testing.T) {
	s := NewStore(store.Prefix("aaa"))
	s.Init(store.Prefix(""))
	if len(s.Options().Prefix) > 0 {
		t.Error("Init didn't reinitialise the store")
	}
}

func TestMemoryBasic(t *testing.T) {
	s := NewStore()
	s.Init()
	basictest(s, t)
}

func TestMemoryPrefix(t *testing.T) {
	s := NewStore()
	s.Init(store.Prefix("some-prefix"))
	basictest(s, t)
}

func TestMemoryNamespace(t *testing.T) {
	s := NewStore()
	s.Init(store.Namespace("some-namespace"))
	basictest(s, t)
}

func TestMemoryNamespacePrefix(t *testing.T) {
	s := NewStore()
	s.Init(store.Prefix("some-prefix"), store.Namespace("some-namespace"))
	basictest(s, t)
}

func basictest(s store.Store, t *testing.T) {
	t.Logf("Testing store %s, with options %# v\n", s.String(), pretty.Formatter(s.Options()))
	// Read and Write an expiring Record
	if err := s.Write(&store.Record{
		Key:    "Hello",
		Value:  []byte("World"),
		Expiry: time.Millisecond * 100,
	}); err != nil {
		t.Error(err)
	}
	if r, err := s.Read("Hello"); err != nil {
		t.Error(err)
	} else {
		if len(r) != 1 {
			t.Error("Read returned multiple records")
		}
		if r[0].Key != "Hello" {
			t.Errorf("Expected %s, got %s", "Hello", r[0].Key)
		}
		if string(r[0].Value) != "World" {
			t.Errorf("Expected %s, got %s", "World", r[0].Value)
		}
	}
	time.Sleep(time.Millisecond * 200)
	if _, err := s.Read("Hello"); err != store.ErrNotFound {
		t.Errorf("Expected %# v, got %# v", store.ErrNotFound, err)
	}

	// Write 3 records with various expiry and get with prefix
	records := []*store.Record{
		&store.Record{
			Key:   "foo",
			Value: []byte("foofoo"),
		},
		&store.Record{
			Key:    "foobar",
			Value:  []byte("foobarfoobar"),
			Expiry: time.Millisecond * 100,
		},
		&store.Record{
			Key:    "foobarbaz",
			Value:  []byte("foobarbazfoobarbaz"),
			Expiry: 2 * time.Millisecond * 100,
		},
	}
	for _, r := range records {
		if err := s.Write(r); err != nil {
			t.Errorf("Couldn't write k: %s, v: %# v (%s)", r.Key, pretty.Formatter(r.Value), err)
		}
	}
	if results, err := s.Read("foo", store.ReadPrefix()); err != nil {
		t.Errorf("Couldn't read all \"foo\" keys, got %# v (%s)", pretty.Formatter(results), err)
	} else {
		if len(results) != 3 {
			t.Errorf("Expected 3 items, got %d", len(results))
		}
		t.Logf("Prefix test: %v\n", pretty.Formatter(results))
	}
	time.Sleep(time.Millisecond * 100)
	if results, err := s.Read("foo", store.ReadPrefix()); err != nil {
		t.Errorf("Couldn't read all \"foo\" keys, got %# v (%s)", pretty.Formatter(results), err)
	} else {
		if len(results) != 2 {
			t.Errorf("Expected 2 items, got %d", len(results))
		}
		t.Logf("Prefix test: %v\n", pretty.Formatter(results))
	}
	time.Sleep(time.Millisecond * 100)
	if results, err := s.Read("foo", store.ReadPrefix()); err != nil {
		t.Errorf("Couldn't read all \"foo\" keys, got %# v (%s)", pretty.Formatter(results), err)
	} else {
		if len(results) != 1 {
			t.Errorf("Expected 1 item, got %d", len(results))
		}
		t.Logf("Prefix test: %# v\n", pretty.Formatter(results))
	}
	if err := s.Delete("foo"); err != nil {
		t.Errorf("Delete failed (%v)", err)
	}
	if results, err := s.Read("foo", store.ReadPrefix()); err != nil {
		t.Errorf("Couldn't read all \"foo\" keys, got %# v (%s)", pretty.Formatter(results), err)
	} else {
		if len(results) != 0 {
			t.Errorf("Expected 0 items, got %d (%# v)", len(results), pretty.Formatter(results))
		}
	}

	// Write 3 records with various expiry and get with Suffix
	records = []*store.Record{
		&store.Record{
			Key:   "foo",
			Value: []byte("foofoo"),
		},
		&store.Record{
			Key:    "barfoo",
			Value:  []byte("barfoobarfoo"),
			Expiry: time.Millisecond * 100,
		},
		&store.Record{
			Key:    "bazbarfoo",
			Value:  []byte("bazbarfoobazbarfoo"),
			Expiry: 2 * time.Millisecond * 100,
		},
	}
	for _, r := range records {
		if err := s.Write(r); err != nil {
			t.Errorf("Couldn't write k: %s, v: %# v (%s)", r.Key, pretty.Formatter(r.Value), err)
		}
	}
	if results, err := s.Read("foo", store.ReadSuffix()); err != nil {
		t.Errorf("Couldn't read all \"foo\" keys, got %# v (%s)", pretty.Formatter(results), err)
	} else {
		if len(results) != 3 {
			t.Errorf("Expected 3 items, got %d", len(results))
		}
		t.Logf("Prefix test: %v\n", pretty.Formatter(results))
	}
	time.Sleep(time.Millisecond * 100)
	if results, err := s.Read("foo", store.ReadSuffix()); err != nil {
		t.Errorf("Couldn't read all \"foo\" keys, got %# v (%s)", pretty.Formatter(results), err)
	} else {
		if len(results) != 2 {
			t.Errorf("Expected 2 items, got %d", len(results))
		}
		t.Logf("Prefix test: %v\n", pretty.Formatter(results))
	}
	time.Sleep(time.Millisecond * 100)
	if results, err := s.Read("foo", store.ReadSuffix()); err != nil {
		t.Errorf("Couldn't read all \"foo\" keys, got %# v (%s)", pretty.Formatter(results), err)
	} else {
		if len(results) != 1 {
			t.Errorf("Expected 1 item, got %d", len(results))
		}
		t.Logf("Prefix test: %# v\n", pretty.Formatter(results))
	}
	if err := s.Delete("foo"); err != nil {
		t.Errorf("Delete failed (%v)", err)
	}
	if results, err := s.Read("foo", store.ReadSuffix()); err != nil {
		t.Errorf("Couldn't read all \"foo\" keys, got %# v (%s)", pretty.Formatter(results), err)
	} else {
		if len(results) != 0 {
			t.Errorf("Expected 0 items, got %d (%# v)", len(results), pretty.Formatter(results))
		}
	}

	// Test Prefix, Suffix and WriteOptions
	if err := s.Write(&store.Record{
		Key:   "foofoobarbar",
		Value: []byte("something"),
	}, store.WriteTTL(time.Millisecond*100)); err != nil {
		t.Error(err)
	}
	if err := s.Write(&store.Record{
		Key:   "foofoo",
		Value: []byte("something"),
	}, store.WriteExpiry(time.Now().Add(time.Millisecond*100))); err != nil {
		t.Error(err)
	}
	if err := s.Write(&store.Record{
		Key:   "barbar",
		Value: []byte("something"),
		// TTL has higher precedence than expiry
	}, store.WriteExpiry(time.Now().Add(time.Hour)), store.WriteTTL(time.Millisecond*100)); err != nil {
		t.Error(err)
	}
	if results, err := s.Read("foo", store.ReadPrefix(), store.ReadSuffix()); err != nil {
		t.Error(err)
	} else {
		if len(results) != 1 {
			t.Errorf("Expected 1 results, got %d: %# v", len(results), pretty.Formatter(results))
		}
	}
	time.Sleep(time.Millisecond * 100)
	if results, err := s.List(); err != nil {
		t.Errorf("List failed: %s", err)
	} else {
		if len(results) != 0 {
			t.Error("Expiry options were not effective")
		}
	}
}
