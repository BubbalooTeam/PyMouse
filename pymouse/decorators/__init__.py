from typing import Optional
from functools import wraps
from hydrogram.types import Message

from pymouse import PyMouse, Config, db, usersmodel_db

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
                if not m.from_user:
                    return await func(c, m, *args, **kwargs)
                user_id = m.from_user.id
                actual_username = m.from_user.username
                actual_uname = m.from_user.first_name + " " + m.from_user.last_name if not m.from_user.last_name is None else m.from_user.first_name
                actual_languagebytg = m.from_user.language_code
                usersmodel_db.users_db.update_user(user_id, actual_uname, actual_username, actual_languagebytg)
                return await func(c, m, *args, **kwargs)
            return wrapper
        return decorator
