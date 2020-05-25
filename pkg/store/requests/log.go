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
	"io"
	"sync"
)

// Index is a log index
type Index uint64

// newLog creates a new in-memory Log
func newLog() Log {
	log := &memoryLog{
		entries:    make([]*Entry, 0, 1024),
		firstIndex: 1,
		readers:    make([]*memoryReader, 0, 10),
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

	// Write appends the given entry to the log
	Write(value interface{}) *Entry
}

// Reader supports reading of entries from the Raft log
type Reader interface {
	io.Closer

	// Index returns the reader index
	Index() Index

	// Read reads the next entry in the log
	Read() *Entry

	// Reset resets the log reader to the given index
	Reset(index Index)
}

// Entry is an indexed Raft log entry
type Entry struct {
	Index Index
	Value interface{}
}

type memoryLog struct {
	entries    []*Entry
	firstIndex Index
	writer     *memoryWriter
	readers    []*memoryReader
	mu         sync.RWMutex
}

func (l *memoryLog) Writer() Writer {
	return l.writer
}

func (l *memoryLog) OpenReader(index Index) Reader {
	readerIndex := -1
	for i := 0; i < len(l.entries); i++ {
		if l.entries[i].Index == index {
			readerIndex = i - 1
			break
		}
	}
	reader := &memoryReader{
		log:   l,
		index: readerIndex,
	}
	l.readers = append(l.readers, reader)
	return reader
}

func (l *memoryLog) Close() error {
	return nil
}

type memoryWriter struct {
	log *memoryLog
}

func (w *memoryWriter) nextIndex() Index {
	if len(w.log.entries) == 0 {
		return w.log.firstIndex
	}
	return w.log.entries[len(w.log.entries)-1].Index + 1
}

func (w *memoryWriter) Write(value interface{}) *Entry {
	w.log.mu.Lock()
	entry := &Entry{
		Index: w.nextIndex(),
		Value: value,
	}
	w.log.entries = append(w.log.entries, entry)
	readers := w.log.readers
	w.log.mu.Unlock()
	for _, reader := range readers {
		reader.next()
	}
	return entry
}

func (w *memoryWriter) Close() error {
	return nil
}

type memoryReader struct {
	log   *memoryLog
	index int
	wg    *sync.WaitGroup
	mu    sync.Mutex
}

func (r *memoryReader) next() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.wg != nil {
		r.wg.Done()
		r.wg = nil
	}
}

func (r *memoryReader) Index() Index {
	if len(r.log.entries) == 0 {
		return r.log.firstIndex - 1
	}
	return r.log.entries[len(r.log.entries)-1].Index
}

func (r *memoryReader) Read() *Entry {
	r.log.mu.RLock()
	if len(r.log.entries) > r.index+1 {
		r.index++
		entry := r.log.entries[r.index]
		r.log.mu.RUnlock()
		return entry
	}
	r.mu.Lock()
	r.wg = &sync.WaitGroup{}
	r.wg.Add(1)
	r.mu.Unlock()
	r.log.mu.RUnlock()
	r.wg.Wait()
	return r.Read()
}

func (r *memoryReader) Reset(index Index) {
	r.log.mu.RLock()
	defer r.log.mu.RUnlock()
	for i := 0; i < len(r.log.entries); i++ {
		if r.log.entries[i].Index >= index {
			r.index = i - 1
			break
		}
	}
}

func (r *memoryReader) Close() error {
	return nil
}
