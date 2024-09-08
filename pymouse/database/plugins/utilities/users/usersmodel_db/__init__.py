#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


from pymouse import db, log

class UsersDB:
    @staticmethod
    def find_user(user_id: int, actual_uname: str, actual_username: str, actual_tglang: str) -> bool:
        users_db = db.GetCollection("users")
        userinfo = users_db.find_one({"user_id": user_id})
        if userinfo:
            username = userinfo.get("username")
            uname = userinfo.get("full_name", "")
            tglang = userinfo.get("tg_lang")
            if username == actual_username and uname == actual_uname and tglang == actual_tglang:
                return True
            return False
        return False
    
    def update_user(self, user_id: int, actual_uname: str, actual_username: str, actual_tglang: str):
        users_db = db.GetCollection("users")
        finders = self.find_user(user_id, actual_uname, actual_username, actual_tglang)
        if not finders:
             users_db.insert_or_update(
                 filter={"user_id": user_id},
                 info={
                     "full_name": actual_uname,
                     "tg_lang": actual_tglang,
                     "username": actual_username,
                     "user_id": user_id,
                 }
             )
             log.info("[database/plugins/utilities/users][UpdateUser]: (%s - %s) User updated with successfully!!", user_id, actual_uname)
             return True
        else:
            return None
        
    @staticmethod
    def getuser_dict(user_id: int) -> dict:
        users_db = db.GetCollection("users")
        return users_db.find_one({"user_id": user_id})

users_db = UsersDB()