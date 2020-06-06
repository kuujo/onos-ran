// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package requests

import (
	"container/list"
	"io"
	"sync"
)

// Index is a log index
type Index uint64

// newLog creates a new in-memory Log
func newLog() Log {
	log := &memoryLog{
		entries: list.New(),
		readers: make([]*memoryReader, 0, 10),
	}
	log.writer = &memoryWriter{
		log: log,
	}
	return log
}

// Log provides for reading and writing entries in the Raft log
type Log interface {
	io.Closer

	// Writer returns the Raft log writer
	Writer() Writer

	// OpenReader opens a Raft log reader
	OpenReader(index Index) Reader
}

// Writer supports writing entries to the Raft log
type Writer interface {
	io.Closer

	// Index returns the writer index
	Index() Index

	// Write appends the given entry to the log
	Write(value interface{}) *Entry

	// Seek seeks to the given index
	Seek(Index)

	// Discard discards the entry at the given index
	Discard(Index)
}

// Reader supports reading of entries from the Raft log
type Reader interface {
	io.Closer

	// ReadBatch reads the next batch from the log
	ReadBatch() *Batch

	// ReadUntil reads the next batch from the log up to the given index
	ReadUntil(Index) *Batch

	// Seek resets the log reader to the given index
	Seek(index Index)
}

// Entry is an indexed Raft log entry
type Entry struct {
	Index Index
	Value interface{}
}

// Batch is a batch of log entries
type Batch struct {
	PrevIndex Index
	Entries   []*Entry
}

type memoryLog struct {
	entries *list.List
	writer  *memoryWriter
	readers []*memoryReader
	mu      sync.RWMutex
}

func (l *memoryLog) Writer() Writer {
	return l.writer
}

func (l *memoryLog) OpenReader(index Index) Reader {
	l.mu.Lock()
	defer l.mu.Unlock()
	var elem *list.Element
	if index > 0 {
		elem = l.entries.Front()
		next := elem
		for next != nil && next.Value.(*Entry).Index <= index {
			elem = next
			next = next.Next()
		}
	}
	reader := &memoryReader{
		log:  l,
		elem: elem,
	}
	l.readers = append(l.readers, reader)
	return reader
}

func (l *memoryLog) Close() error {
	return nil
}

type memoryWriter struct {
	log       *memoryLog
	lastIndex Index
}

func (w *memoryWriter) Index() Index {
	w.log.mu.RLock()
	defer w.log.mu.RUnlock()
	return w.lastIndex
}

func (w *memoryWriter) nextIndex() Index {
	w.lastIndex++
	return w.lastIndex
}

func (w *memoryWriter) Write(value interface{}) *Entry {
	w.log.mu.Lock()
	index := w.nextIndex()
	entry := &Entry{
		Index: index,
		Value: value,
	}
	w.log.entries.PushBack(entry)
	for _, reader := range w.log.readers {
		reader.next()
	}
	w.log.mu.Unlock()
	return entry
}

func (w *memoryWriter) Seek(index Index) {
	w.log.mu.Lock()
	defer w.log.mu.Unlock()
	w.lastIndex = index
}

func (w *memoryWriter) Discard(index Index) {
	w.log.mu.Lock()
	defer w.log.mu.Unlock()
	elem := w.log.entries.Front()
	for elem != nil && elem.Value.(*Entry).Index <= index {
		w.log.entries.Remove(elem)
		elem = elem.Next()
	}
}

func (w *memoryWriter) Close() error {
	return nil
}

type memoryReader struct {
	log     *memoryLog
	elem    *list.Element
	wg      *sync.WaitGroup
	batchCh chan *Batch
	mu      sync.Mutex
}

func (r *memoryReader) next() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.wg != nil {
		r.wg.Done()
		r.wg = nil
	}
}

func (r *memoryReader) lock() {
	r.log.mu.RLock()
	r.mu.Lock()
}

func (r *memoryReader) unlock() {
	r.mu.Unlock()
	r.log.mu.RUnlock()
}

func (r *memoryReader) waiter() *sync.WaitGroup {
	if r.wg == nil {
		r.wg = &sync.WaitGroup{}
		r.wg.Add(1)
	}
	return r.wg
}

func (r *memoryReader) ReadBatch() *Batch {
	r.lock()
	if r.elem != nil && (r.log.entries.Len() == 0 || r.elem.Value.(*Entry).Index < r.log.entries.Front().Value.(*Entry).Index) {
		r.elem = nil
	}
	if (r.elem == nil && r.log.entries.Len() > 0) || (r.elem != nil && r.elem.Next() != nil) {
		var firstElem *list.Element
		if r.elem == nil {
			firstElem = r.log.entries.Front()
		} else {
			firstElem = r.elem.Next()
		}

		entries := make([]*Entry, 0)
		elem := firstElem
		for elem != nil {
			entries = append(entries, elem.Value.(*Entry))
			r.elem = elem
			elem = elem.Next()
		}

		var prevIndex Index
		if firstElem != nil {
			prevIndex = firstElem.Value.(*Entry).Index - 1
		}
		r.mu.Unlock()
		r.log.mu.RUnlock()
		return &Batch{
			PrevIndex: prevIndex,
			Entries:   entries,
		}
	}
	wg := r.waiter()
	r.unlock()
	wg.Wait()
	return r.ReadBatch()
}

func (r *memoryReader) ReadUntil(index Index) *Batch {
	r.lock()
	if r.elem != nil && (r.log.entries.Len() == 0 || r.elem.Value.(*Entry).Index < r.log.entries.Front().Value.(*Entry).Index) {
		r.elem = nil
	}
	if r.log.writer.lastIndex >= index && (r.elem == nil && r.log.entries.Len() > 0) || (r.elem != nil && r.elem.Next() != nil) {
		var firstElem *list.Element
		if r.elem == nil {
			firstElem = r.log.entries.Front()
		} else {
			firstElem = r.elem.Next()
		}

		entries := make([]*Entry, 0)
		elem := firstElem
		for elem != nil && elem.Value.(*Entry).Index <= index {
			entries = append(entries, elem.Value.(*Entry))
			r.elem = elem
			elem = elem.Next()
		}

		var prevIndex Index
		if firstElem != nil {
			prevIndex = firstElem.Value.(*Entry).Index - 1
		}
		r.mu.Unlock()
		r.log.mu.RUnlock()
		return &Batch{
			PrevIndex: prevIndex,
			Entries:   entries,
		}
	}
	wg := r.waiter()
	r.unlock()
	wg.Wait()
	return r.ReadUntil(index)
}

func (r *memoryReader) Seek(index Index) {
	r.lock()
	defer r.unlock()
	if r.elem == nil {
		var elem *list.Element
		next := r.log.entries.Front()
		for next != nil && next.Value.(*Entry).Index < index-1 {
			elem = next
			next = elem.Next()
		}
		r.elem = elem
	} else if r.elem.Value.(*Entry).Index > index-1 {
		elem := r.elem
		prev := elem.Prev()
		for prev != nil && prev.Value.(*Entry).Index >= index - 1 {
			elem = prev
			prev = elem.Prev()
		}
		if elem.Value.(*Entry).Index > index-1 {
			r.elem = nil
		} else {
			r.elem = elem
		}
	} else if r.elem.Value.(*Entry).Index < index-1 {
		elem := r.elem
		next := elem.Next()
		for next != nil && next.Value.(*Entry).Index < index-1 {
			elem = next
			next = elem.Next()
		}
		r.elem = elem
	}
}

func (r *memoryReader) Close() error {
	return nil
}
