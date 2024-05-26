from pymouse import db, log

class ChatsDB:
    @staticmethod
    def find_chat(chat_id: int, actual_title: str, actual_chatname: str) -> bool:
        chats_db = db.GetCollection("chats")
        chatinfo = chats_db.find_one({"chat_id": chat_id})
        if chatinfo:
            chatname = chatinfo.get("chatname", None)
            chattitle = chatinfo.get("chattitle", "")
            if chatname == actual_chatname and chattitle == actual_title:
                return True
            return False
        return False

    def update_chat(self, chat_id: int, actual_chattitle: str, actual_chatname: str):
        chats_db = db.GetCollection("chats")
        finders = self.find_chat(chat_id, actual_chattitle, actual_chatname)
        if not finders:
            chats_db.insert_or_update(
                filter={"chat_id": chat_id},
                info={
                    "chattitle": actual_chattitle,
                    "chatname": actual_chatname,
                    "chat_id": chat_id
                }
            )
            log.info("[database/plugins/utilities/chats]: (%s - %s) Chat updated with successfully!!", chat_id, actual_chattitle)
            return True
        else:
            return None
        
    @staticmethod
    def get_chat_dict(chat_id: int) -> dict:
        chats_db = db.GetCollection("chats")
        return chats_db.find_one({"chat_id": chat_id})


chats_db = ChatsDB()