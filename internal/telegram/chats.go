package telegram

import (
	"context"
	"fmt"

	"github.com/gotd/td/constant"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

// LoadChats fetches the dialog list and pushes every chat to the UI
// as a ChatUpdateMsg (this replaces tdlib's updateNewChat flow).
func (c *Client) LoadChats(limit int) error {
	chats, err := c.ListChats(limit)
	if err != nil {
		return err
	}

	var totalUnread int32
	for _, chat := range chats {
		totalUnread += chat.UnreadCount
		c.send(ChatUpdateMsg{Chat: chat})
	}

	c.send(UnreadCountMsg{
		UnreadCount:        totalUnread,
		UnreadUnmutedCount: totalUnread,
	})
	return nil
}

// ListChats fetches the dialog list without emitting any UI events.
func (c *Client) ListChats(limit int) ([]*Chat, error) {
	ctx := context.Background()
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	res, err := c.api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		Limit:      limit,
		OffsetPeer: &tg.InputPeerEmpty{},
	})
	if err != nil {
		return nil, fmt.Errorf("get dialogs: %w", err)
	}

	var (
		dialogs  []tg.DialogClass
		messages []tg.MessageClass
		chats    []tg.ChatClass
		users    []tg.UserClass
	)
	switch d := res.(type) {
	case *tg.MessagesDialogs:
		dialogs, messages, chats, users = d.Dialogs, d.Messages, d.Chats, d.Users
	case *tg.MessagesDialogsSlice:
		dialogs, messages, chats, users = d.Dialogs, d.Messages, d.Chats, d.Users
	default:
		return nil, fmt.Errorf("unexpected dialogs type %T", res)
	}

	// Seed the peers manager so access hashes are known.
	if err := c.peers.Apply(ctx, users, chats); err != nil {
		return nil, fmt.Errorf("apply peers: %w", err)
	}

	entities := tg.Entities{
		Users:    make(map[int64]*tg.User, len(users)),
		Chats:    make(map[int64]*tg.Chat, len(chats)),
		Channels: make(map[int64]*tg.Channel, len(chats)),
	}
	for _, uc := range users {
		if u, ok := uc.(*tg.User); ok {
			entities.Users[u.ID] = u
		}
	}
	for _, cc := range chats {
		switch v := cc.(type) {
		case *tg.Chat:
			entities.Chats[v.ID] = v
		case *tg.Channel:
			entities.Channels[v.ID] = v
		}
	}

	lastMessages := make(map[int64]*Message, len(messages))
	for _, mc := range messages {
		if m := c.messageClassFromTG(mc); m != nil {
			lastMessages[m.ChatID] = m
		}
	}

	out := make([]*Chat, 0, len(dialogs))
	for _, dc := range dialogs {
		d, ok := dc.(*tg.Dialog)
		if !ok {
			continue
		}

		chat, err := c.chatFromPeer(d.Peer, entities)
		if err != nil {
			continue
		}

		chat.Pinned = d.Pinned
		chat.UnreadCount = int32(d.UnreadCount)
		chat.LastReadInboxMessageID = int64(d.ReadInboxMaxID)
		chat.LastReadOutboxMessageID = int64(d.ReadOutboxMaxID)
		if lm, ok := lastMessages[chat.ID]; ok {
			chat.LastMessage = lm
			chat.Order = int64(lm.Date)
		}

		out = append(out, chat)
	}
	return out, nil
}

// GetChat returns a single chat by canonical chat ID.
func (c *Client) GetChat(chatID int64) (*Chat, error) {
	ctx := context.Background()
	peer, err := c.peers.ResolveTDLibID(ctx, constant.TDLibPeerID(chatID))
	if err != nil {
		return nil, fmt.Errorf("get chat %d: %w", chatID, err)
	}

	switch p := peer.(type) {
	case peers.User:
		return c.chatFromUser(p.Raw()), nil
	case peers.Chat:
		return c.chatFromBasicGroup(p.Raw()), nil
	case peers.Channel:
		return c.chatFromChannel(p.Raw()), nil
	default:
		return nil, fmt.Errorf("get chat %d: unexpected peer type %T", chatID, peer)
	}
}

// GetChatHistory returns messages of a chat, newest first.
// fromMessageID paginates backwards (offsetID); offset skips messages.
func (c *Client) GetChatHistory(chatID, fromMessageID int64, offset, limit int32) ([]*Message, error) {
	ctx := context.Background()
	peer, err := c.inputPeer(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	res, err := c.api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer:       peer,
		OffsetID:   int(fromMessageID),
		OffsetDate: 0,
		AddOffset:  int(offset),
		Limit:      int(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}

	var messages []tg.MessageClass
	switch r := res.(type) {
	case *tg.MessagesMessages:
		messages = r.Messages
	case *tg.MessagesMessagesSlice:
		messages = r.Messages
	case *tg.MessagesChannelMessages:
		messages = r.Messages
	default:
		return nil, fmt.Errorf("unexpected history type %T", res)
	}

	out := make([]*Message, 0, len(messages))
	for _, mc := range messages {
		if m := c.messageClassFromTG(mc); m != nil {
			out = append(out, m)
		}
	}
	return out, nil
}

// SearchChats searches chat titles by query (server-side).
func (c *Client) SearchChats(query string, limit int32) ([]*Chat, error) {
	ctx := context.Background()
	if limit <= 0 {
		limit = 20
	}

	res, err := c.api.ContactsSearch(ctx, &tg.ContactsSearchRequest{
		Q:     query,
		Limit: int(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("search chats: %w", err)
	}

	found := res

	entities := tg.Entities{
		Users:    make(map[int64]*tg.User),
		Chats:    make(map[int64]*tg.Chat),
		Channels: make(map[int64]*tg.Channel),
	}
	for _, uc := range found.Users {
		if u, ok := uc.(*tg.User); ok {
			entities.Users[u.ID] = u
		}
	}
	for _, cc := range found.Chats {
		switch v := cc.(type) {
		case *tg.Chat:
			entities.Chats[v.ID] = v
		case *tg.Channel:
			entities.Channels[v.ID] = v
		}
	}

	out := make([]*Chat, 0, len(found.MyResults)+len(found.Results))
	for _, peer := range append(found.MyResults, found.Results...) {
		if chat, err := c.chatFromPeer(peer, entities); err == nil {
			out = append(out, chat)
		}
	}
	return out, nil
}

// SearchMessages searches messages globally by query.
func (c *Client) SearchMessages(query string, limit int32) ([]*Message, error) {
	ctx := context.Background()
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	res, err := c.api.MessagesSearchGlobal(ctx, &tg.MessagesSearchGlobalRequest{
		Q:          query,
		Filter:     &tg.InputMessagesFilterEmpty{},
		OffsetPeer: &tg.InputPeerEmpty{},
		Limit:      int(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("search messages: %w", err)
	}

	var messages []tg.MessageClass
	switch r := res.(type) {
	case *tg.MessagesMessages:
		messages = r.Messages
	case *tg.MessagesMessagesSlice:
		messages = r.Messages
	case *tg.MessagesChannelMessages:
		messages = r.Messages
	default:
		return nil, fmt.Errorf("unexpected search type %T", res)
	}

	out := make([]*Message, 0, len(messages))
	for _, mc := range messages {
		if m := c.messageClassFromTG(mc); m != nil {
			out = append(out, m)
		}
	}
	return out, nil
}

// OpenChat is a light-weight placeholder kept for API compatibility:
// gotd needs no open/close chat lifecycle. It emits the chat so the
// store has it even for chats outside the loaded dialogs.
func (c *Client) OpenChat(chatID int64) error {
	chat, err := c.GetChat(chatID)
	if err != nil {
		return err
	}
	c.send(ChatUpdateMsg{Chat: chat})
	return nil
}

// ViewMessages marks messages as read.
func (c *Client) ViewMessages(chatID int64, messageIDs []int64) error {
	ctx := context.Background()
	peer, err := c.inputPeer(ctx, chatID)
	if err != nil {
		return fmt.Errorf("view messages: %w", err)
	}

	maxID := int64(0)
	for _, id := range messageIDs {
		if id > maxID {
			maxID = id
		}
	}
	if maxID == 0 {
		return nil
	}

	if constant.TDLibPeerID(chatID).IsChannel() {
		inputChannel, ok := peerAsInputChannel(peer)
		if !ok {
			return fmt.Errorf("view messages: peer %d is not a channel", chatID)
		}
		if _, err := c.api.ChannelsReadHistory(ctx, &tg.ChannelsReadHistoryRequest{
			Channel: inputChannel,
			MaxID:   int(maxID),
		}); err != nil {
			return fmt.Errorf("read channel history: %w", err)
		}
		return nil
	}

	if _, err := c.api.MessagesReadHistory(ctx, &tg.MessagesReadHistoryRequest{
		Peer:  peer,
		MaxID: int(maxID),
	}); err != nil {
		return fmt.Errorf("read history: %w", err)
	}
	return nil
}

// peerAsInputChannel extracts an InputChannel from an InputPeer.
func peerAsInputChannel(peer tg.InputPeerClass) (tg.InputChannelClass, bool) {
	switch p := peer.(type) {
	case *tg.InputPeerChannel:
		return &tg.InputChannel{ChannelID: p.ChannelID, AccessHash: p.AccessHash}, true
	case *tg.InputPeerChannelFromMessage:
		return &tg.InputChannelFromMessage{
			Peer:      p.Peer,
			MsgID:     p.MsgID,
			ChannelID: p.ChannelID,
		}, true
	default:
		return nil, false
	}
}
