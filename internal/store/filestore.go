package store

import (
	"sync"

	"github.com/zelenin/go-tdlib/client"
)

// FileState tracks the download/upload state of a file.
type FileState struct {
	File       *client.File
	LocalPath  string
	IsComplete bool
	Progress   float64 // 0.0 to 1.0
}

// FileStore tracks file download/upload states.
type FileStore struct {
	mu    sync.RWMutex
	files map[int32]*FileState // fileID -> state
}

func NewFileStore() *FileStore {
	return &FileStore{
		files: make(map[int32]*FileState),
	}
}

// Update processes a file update from TDLib.
func (s *FileStore) Update(file *client.File) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := &FileState{
		File: file,
	}

	if file.Local != nil {
		state.LocalPath = file.Local.Path
		state.IsComplete = file.Local.IsDownloadingCompleted

		if file.ExpectedSize > 0 {
			state.Progress = float64(file.Local.DownloadedSize) / float64(file.ExpectedSize)
		}
	}

	s.files[file.ID] = state
}

// Get returns the state of a file.
func (s *FileStore) Get(fileID int32) (*FileState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.files[fileID]
	return state, ok
}

// IsComplete checks if a file download is complete.
func (s *FileStore) IsComplete(fileID int32) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.files[fileID]
	if !ok {
		return false
	}
	return state.IsComplete
}

// LocalPath returns the local path of a downloaded file.
func (s *FileStore) LocalPath(fileID int32) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.files[fileID]
	if !ok {
		return ""
	}
	return state.LocalPath
}

// Progress returns the download progress of a file (0.0 to 1.0).
func (s *FileStore) Progress(fileID int32) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.files[fileID]
	if !ok {
		return 0
	}
	return state.Progress
}

// Store is the aggregate store holding all caches.
type Store struct {
	Chats    *ChatStore
	Messages *MessageStore
	Users    *UserStore
	Files    *FileStore
}

// NewStore creates a new aggregate store.
func NewStore() *Store {
	return &Store{
		Chats:    NewChatStore(),
		Messages: NewMessageStore(),
		Users:    NewUserStore(),
		Files:    NewFileStore(),
	}
}
