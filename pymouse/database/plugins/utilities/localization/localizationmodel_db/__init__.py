from typing import Union
from hydrogram.types import Message, CallbackQuery, InlineQuery
from hydrogram.enums import ChatType

from pymouse import db, log

class LocalizationDB:
    @staticmethod
    def get_chatid(
        union: Union[Message, CallbackQuery, InlineQuery],
    ):
        """
        Returns the chatID.

        Arguments:
            union (Union[Message, CallbackQuery, InlineQuery]): The received message or query.

        Returns:
            int | None: The chat ID or None if not supported.
        """
        if isinstance(union, InlineQuery):
            return union.from_user.id
        elif isinstance(union, Union[Message, CallbackQuery]):
            msg = union.message if isinstance(union, CallbackQuery) else union
            return msg.chat.id
        else:
            log.critical("This instance is not supported!")
            return None

    @staticmethod
    def get_database_context(
        union: Union[Message, CallbackQuery, InlineQuery],
    ):
        """
        Returns the appropriate database collection.

        Arguments:
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
            elif chat_type in {ChatType.PRIVATE, ChatType.BOT}:
                return db.GetCollection("users")
            else:
                log.critical("Chat type with enum (%s) is not supported!", chat_type)
                return None
        else:
            log.critical("This instance is not supported!")
            return None
            
    def get_chat_language(
        self,
        union: Union[Message, CallbackQuery, InlineQuery],
    ):
        """
        Returns the chat language.

        Arguments:
            union (Union[Message, CallbackQuery, InlineQuery]): The received message or query.

        Returns:
            str | None: The chat language or None if not supported.
        """
        db_ctx = self.get_database_context(union)
        chat_id = self.get_chatid(union)
        tck_type = "chat_id" if db_ctx.collection == "chats" else "user_id"
        get_chat = db_ctx.find_one({tck_type: chat_id})
        if get_chat is not None:
            language_class = get_chat.get("localization", {})
            language = language_class.get("chat_lang", None)
            return language
        else:
            return None
        
    def set_chat_language(
        self,
        union: Union[Message, CallbackQuery],
        language: str,
    ):
        db_ctx = self.get_database_context(union)
        chat_id = self.get_chatid(union)
        tck_type = "chat_id" if db_ctx.collection == "chats" else "user_id"
        language_mapper = {
            tck_type: chat_id,
            "localization": {
                "chat_lang": language,
            }
        }
        db_ctx.insert_or_update(
            filter={tck_type: chat_id},
            info=language_mapper,
        )


localization_db = LocalizationDB()