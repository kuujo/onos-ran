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

import "io"

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

	// LastIndex returns the last index written to the log
	LastIndex() Index

	// LastEntry returns the last entry written to the log
	LastEntry() *Entry

	// Append appends the given entry to the log
	Append(request *Request) *Entry

	// Reset resets the log writer to the given index
	Reset(index Index)

	// Truncate truncates the tail of the log to the given index
	Truncate(index Index)
}

// Reader supports reading of entries from the Raft log
type Reader interface {
	io.Closer

	// FirstIndex returns the first index in the log
	FirstIndex() Index

	// LastIndex returns the last index in the log
	LastIndex() Index

	// CurrentIndex returns the current index of the reader
	CurrentIndex() Index

	// CurrentEntry returns the current Entry
	CurrentEntry() *Entry

	// NextIndex returns the next index in the log
	NextIndex() Index

	// NextEntry advances the log index and returns the next entry in the log
	NextEntry() *Entry

	// Reset resets the log reader to the given index
	Reset(index Index)
}

// Entry is an indexed Raft log entry
type Entry struct {
	Index   Index
	Request *Request
}

type memoryLog struct {
	entries    []*Entry
	firstIndex Index
	writer     *memoryWriter
	readers    []*memoryReader
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

func (w *memoryWriter) LastIndex() Index {
	if entry := w.LastEntry(); entry != nil {
		return entry.Index
	}
	return w.log.firstIndex - 1
}

func (w *memoryWriter) LastEntry() *Entry {
	if len(w.log.entries) == 0 {
		return nil
	}
	return w.log.entries[len(w.log.entries)-1]
}

func (w *memoryWriter) nextIndex() Index {
	if len(w.log.entries) == 0 {
		return w.log.firstIndex
	}
	return w.log.entries[len(w.log.entries)-1].Index + 1
}

func (w *memoryWriter) Append(request *Request) *Entry {
	indexed := &Entry{
		Index:   w.nextIndex(),
		Request: request,
	}
	w.log.entries = append(w.log.entries, indexed)
	return indexed
}

func (w *memoryWriter) Reset(index Index) {
	w.log.entries = w.log.entries[:0]
	w.log.firstIndex = index
	for _, reader := range w.log.readers {
		reader.maybeReset()
	}
}

func (w *memoryWriter) Truncate(index Index) {
	for i := 0; i < len(w.log.entries); i++ {
		if w.log.entries[i].Index > index {
			w.log.entries = w.log.entries[:i]
			break
		}
	}
	for _, reader := range w.log.readers {
		reader.maybeReset()
	}
}

func (w *memoryWriter) Close() error {
	panic("implement me")
}

type memoryReader struct {
	log   *memoryLog
	index int
}

func (r *memoryReader) FirstIndex() Index {
	return r.log.firstIndex
}

func (r *memoryReader) LastIndex() Index {
	if len(r.log.entries) == 0 {
		return r.log.firstIndex - 1
	}
	return r.log.entries[len(r.log.entries)-1].Index
}

func (r *memoryReader) CurrentIndex() Index {
	if r.index == -1 || len(r.log.entries) == 0 {
		return r.log.firstIndex - 1
	}
	return r.log.entries[r.index].Index
}

func (r *memoryReader) CurrentEntry() *Entry {
	if r.index == -1 || len(r.log.entries) == 0 {
		return nil
	}
	return r.log.entries[r.index]
}

func (r *memoryReader) NextIndex() Index {
	if r.index == -1 || len(r.log.entries) == 0 {
		return r.log.firstIndex
	}
	return r.log.entries[r.index].Index + 1
}

func (r *memoryReader) NextEntry() *Entry {
	if len(r.log.entries) > r.index+1 {
		r.index++
		return r.log.entries[r.index]
	}
	return nil
}

func (r *memoryReader) Reset(index Index) {
	for i := 0; i < len(r.log.entries); i++ {
		if r.log.entries[i].Index >= index {
			r.index = i - 1
			break
		}
	}
}

func (r *memoryReader) maybeReset() {
	if r.index >= 0 && len(r.log.entries) <= r.index {
		r.index = len(r.log.entries) - 1
	}
}

func (r *memoryReader) Close() error {
	return nil
}
