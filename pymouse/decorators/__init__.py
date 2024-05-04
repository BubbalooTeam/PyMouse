from typing import Optional
from functools import wraps
from hydrogram.types import Message

from pymouse import PyMouse, Config, db

class Decorators:
    def __init__(self):
        self.dev_db = db.GetCollection("SUDOERS")
        self.owner = Config.OWNER_ID
    
    def require_dev(self, only_owner: Optional[bool] = False):
        def decorator(func):
            @wraps(func)
            async def wrapper(c: PyMouse, m: Message, *args, **kwargs):
                if not m.from_user:
                    return None
                user_id = m.from_user.id
                if user_id == self.owner:
                    return await func(c, m, *args, **kwargs)
                get_dev_by_id = self.dev_db.find_one({"dev_users": user_id})
                if get_dev_by_id and not only_owner is True:
                    return await func(c, m, *args, **kwargs)
                else:
                    return None
            return wrapper
        return decorator