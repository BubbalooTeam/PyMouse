from hydrogram.types import Message
from pymouse import PyMouse, Decorators

@PyMouse.on_message(group=1)
@Decorators().SaveUsers()
@Decorators().SaveChats()
async def bot_utilities(c: PyMouse, m: Message): # type: ignore
    pass