package store

import (
	"sort"
	"sync"

	"github.com/zelenin/go-tdlib/client"
)

// ChatEntry holds chat metadata and its position in the main list.
type ChatEntry struct {
	Chat        *client.Chat
	LastMessage *client.Message
	UnreadCount int32
	Position    int64
}

// ChatStore is a thread-safe in-memory cache of chats.
type ChatStore struct {
	mu    sync.RWMutex
	chats map[int64]*ChatEntry
}

func NewChatStore() *ChatStore {
	return &ChatStore{
		chats: make(map[int64]*ChatEntry),
	}
}

// Set adds or updates a chat entry.
func (s *ChatStore) Set(chat *client.Chat) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.chats[chat.ID]
	if !exists {
		entry = &ChatEntry{}
		s.chats[chat.ID] = entry
	}
	entry.Chat = chat
	entry.UnreadCount = chat.UnreadCount

	if chat.LastMessage != nil {
		entry.LastMessage = chat.LastMessage
	}

	if len(chat.Positions) > 0 {
		for _, pos := range chat.Positions {
			if _, ok := pos.List.(*client.ChatListMain); ok {
				entry.Position = pos.Order
			}
		}
	}
}

// Get returns a chat entry by ID.
func (s *ChatStore) Get(chatID int64) (*ChatEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.chats[chatID]
	return entry, ok
}

// UpdatePosition updates a chat's sort position.
func (s *ChatStore) UpdatePosition(chatID int64, positions []*client.ChatPosition) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.chats[chatID]
	if !ok {
		entry = &ChatEntry{Chat: &client.Chat{ID: chatID}}
		s.chats[chatID] = entry
	}

	for _, pos := range positions {
		if _, ok := pos.List.(*client.ChatListMain); ok {
			entry.Position = pos.Order
		}
	}
}

// UpdateLastMessage updates a chat's last message and position.
func (s *ChatStore) UpdateLastMessage(chatID int64, msg *client.Message, positions []*client.ChatPosition) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.chats[chatID]
	if !ok {
		entry = &ChatEntry{Chat: &client.Chat{ID: chatID}}
		s.chats[chatID] = entry
	}

	entry.LastMessage = msg
	for _, pos := range positions {
		if _, ok := pos.List.(*client.ChatListMain); ok {
			entry.Position = pos.Order
		}
	}
}

// UpdateReadInbox updates the unread count for a chat.
func (s *ChatStore) UpdateReadInbox(chatID int64, unreadCount int32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, ok := s.chats[chatID]; ok {
		entry.UnreadCount = unreadCount
	}
}

// OrderedChats returns all chats sorted by position (descending).
func (s *ChatStore) OrderedChats() []*ChatEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries := make([]*ChatEntry, 0, len(s.chats))
	for _, entry := range s.chats {
		if entry.Position > 0 {
			entries = append(entries, entry)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Position > entries[j].Position
	})

	return entries
}

// Count returns the number of cached chats.
func (s *ChatStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.chats)
}
