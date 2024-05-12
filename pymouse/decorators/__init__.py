from typing import Optional
from functools import wraps
from hydrogram.types import Message

from pymouse import PyMouse, Config, db, log

class Decorators:
    def __init__(self):
        self.owner = Config.OWNER_ID

    def require_dev(self, only_owner: Optional[bool] = False):
        def decorator(func):
            @wraps(func)
            async def wrapper(c: PyMouse, m: Message, *args, **kwargs): # type: ignore
                dev_db = db.GetCollection("sudoers")
                if not m.from_user:
                    return None
                user_id = m.from_user.id
                if user_id == self.owner:
                    return await func(c, m, *args, **kwargs)
                get_dev_by_id = dev_db.find_one({"dev_users": user_id})
                if get_dev_by_id and only_owner is False:
                    return await func(c, m, *args, **kwargs)
                else:
                    return None
            return wrapper
        return decorator
    
    def SaveUsers(self):
        def decorator(func):
            @wraps(func)
            async def wrapper(c: PyMouse, m: Message, *args, **kwargs): # type: ignore
                users_db = db.GetCollection("users")
                if not m.from_user:
                    return await func(c, m, *args, **kwargs)
                user_id = m.from_user.id
                actual_username = m.from_user.username
                actual_uname = m.from_user.first_name + " " + m.from_user.last_name if not m.from_user.last_name is None else m.from_user.first_name
                actual_languagebytg = m.from_user.language_code
                # Check if user in this db, if not in DB try to updates
                get_users = users_db.find_one({"user_id": user_id})
                if get_users:
                    # Check Information's of this user
                    username = get_users.get("username", "")
                    uname = get_users.get("full_name", "")
                    tglang = get_users.get("tglang", "N/A")
                    if username == actual_username and uname == actual_uname and tglang == actual_languagebytg:
                        return None
                    users_db.insert_or_update(
                        filter={"user_id": user_id},
                        info={
                            "full_name": actual_uname,
                            "username": actual_username,
                            "tglang": actual_languagebytg,
                        }
                    )
                    return await func(c, m, *args, *kwargs)
                users_db.insert_or_update(
                    filter={"user_id": user_id},
                    info={
                        "full_name": actual_uname,
                        "username": actual_username,
                        "user_id": user_id,
                        "tglang": actual_languagebytg,
                    }
                )
                log.info("[%s | %s]: Added user to LocalDatabase with sucessfully!", user_id, uname)
                return await func(c, m, *args, **kwargs)
            return wrapper
        return decorator
