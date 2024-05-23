from pymouse import db, log

class UsersDB:
    async def find_user(self, user_id: int, actual_uname: str, actual_username: str, actual_tglang: str) -> bool:
        users_db = db.GetCollection("users")
        userinfo = (await users_db.find_one({"user_id": user_id}))
        if userinfo:
            username = userinfo.get("username", None)
            uname = userinfo.get("full_name", "")
            tglang = userinfo.get("tg_lang", None)
            if username == actual_username and uname == actual_uname and tglang == actual_tglang:
                return True
            return False
        return False
    
    async def update_user(self, user_id: int, actual_uname: str, actual_username: str, actual_tglang: str):
        users_db = db.GetCollection("users")
        finders = await self.find_user(user_id, actual_uname, actual_username, actual_tglang)
        if not finders:
             await users_db.insert_or_update(
                 filter={"user_id": user_id},
                 info={
                     "full_name": actual_uname,
                     "tg_lang": actual_tglang,
                     "username": actual_username,
                     "user_id": user_id,
                 }
             )
             log.info("[database/plugins/utilities/users][UpdateUser]: (%s - %s) User updated with sucessfully!!", user_id, actual_uname)
             return True
        else:
            return None
        
    async def getuser_dict(self, user_id: int):
        users_db = db.GetCollection("users")
        return await users_db.find_one({"user_id": user_id})

users_db = UsersDB()