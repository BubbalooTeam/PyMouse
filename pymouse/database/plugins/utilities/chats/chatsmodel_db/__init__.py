from pymouse import db, log

class ChatsDB:
    def find_chat(self, chat_id: int, actual_title: str, actual_chatname: str) -> bool:
        chats_db = db.GetCollection("chats")
        chatinfo = chats_db.find_one({"chat_id": chat_id})
        if chatinfo:
            chatname = chatinfo.get("chatname", None)
            chattitle = chatinfo.get("chattitle", "")
            if chatname == actual_chatname and chattitle == actual_title:
                return True
            return False
        return False
            