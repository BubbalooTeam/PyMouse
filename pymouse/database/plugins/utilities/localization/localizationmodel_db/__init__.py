from typing import Union
from hydrogram.types import Message, CallbackQuery, InlineQuery
from hydrogram.enums import ChatType

from pymouse import db, log

class LocalizationDB:
    def get_database_context(
        self,
        union: Union[Message, CallbackQuery, InlineQuery]
    ):
        """
        Returns the appropriate database collection based on the message or query type.

        Args:
            union (Union[Message, CallbackQuery, InlineQuery]): The received message or query.

        Returns:
            db.Collection: The corresponding database collection.
        """
        if isinstance(union, InlineQuery):
            return db.GetCollection("users")
        elif isinstance(union, (Message, CallbackQuery)):
            msg = union.message if isinstance(union, CallbackQuery) else union
            chat_type = msg.chat.type
            if chat_type in {ChatType.GROUP, ChatType.SUPERGROUP, ChatType.CHANNEL}:
                return db.GetCollection("chats")
            elif chat_type == ChatType.PRIVATE:
                return db.GetCollection("users")
            else:
                log.critical("Chat type with enum (%s) is not supported!", chat_type)
                return None


localization_db = LocalizationDB()