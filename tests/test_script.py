import unittest
from unittest.mock import MagicMock
import contextlib

from internal.websocket import publisher
from internal.dto import Message

class TestChatStore(unittest.TestCase):
    def setUp(self):
        self.chat_store = MagicMock(spec=publisher.ChatStore)
        self.ctx = contextlib.nullcontext()
        self.msg = Message()

    def test_publish_message(self):
        self.chat_store.PublishMessage.return_value = None
        result = self.chat_store.PublishMessage(self.ctx, 'room1', self.msg)
        self.chat_store.PublishMessage.assert_called_with(self.ctx, 'room1', self.msg)
        self.assertIsNone(result)

    def test_save_message(self):
        self.chat_store.SaveMessage.return_value = None
        result = self.chat_store.SaveMessage(self.ctx, 'room1', self.msg, 10)
        self.chat_store.SaveMessage.assert_called_with(self.ctx, 'room1', self.msg, 10)
        self.assertIsNone(result)

    def test_get_messages(self):
        self.chat_store.GetMessages.return_value = ([self.msg], None)
        messages, err = self.chat_store.GetMessages(self.ctx, 'room1', 10)
        self.chat_store.GetMessages.assert_called_with(self.ctx, 'room1', 10)
        self.assertEqual(messages, [self.msg])
        self.assertIsNone(err)

    def test_save_unread(self):
        self.chat_store.SaveUnread.return_value = None
        result = self.chat_store.SaveUnread(self.ctx, 'user1', self.msg)
        self.chat_store.SaveUnread.assert_called_with(self.ctx, 'user1', self.msg)
        self.assertIsNone(result)

    def test_get_unread_messages(self):
        self.chat_store.GetUnreadMessages.return_value = ([self.msg], None)
        messages, err = self.chat_store.GetUnreadMessages(self.ctx, 'user1')
        self.chat_store.GetUnreadMessages.assert_called_with(self.ctx, 'user1')
        self.assertEqual(messages, [self.msg])
        self.assertIsNone(err)

    def test_clear_unread(self):
        self.chat_store.ClearUnread.return_value = None
        result = self.chat_store.ClearUnread(self.ctx, 'user1')
        self.chat_store.ClearUnread.assert_called_with(self.ctx, 'user1')
        self.assertIsNone(result)

if __name__ == '__main__':
    unittest.main()
