#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


from typing import Union

from hydrogram.types import Message, InlineQuery

from typing import Callable
from functools import partial, wraps
from asyncio import get_event_loop

class UtilsDecorator:
    def aiowrap(self):
        """Run syncronomous processes in 2nd plan..."""
        def decorator(func: Callable) -> Callable:
            @wraps(func)
            async def wrapper(*args, loop=None, executor=None, **kwargs):
                if loop is None:
                    loop = get_event_loop()
                pfunc = partial(func, *args, **kwargs)
                return await loop.run_in_executor(executor, pfunc)
            return wrapper
        return decorator
    
class UtilsTimer:
    @staticmethod
    def time_formatter(seconds: float) -> str:
        """
        Show humanized timer.
        Examples: 1d, 1h, 10m, 22s
        """
        minutes, seconds = divmod(int(seconds), 60)
        hours, minutes = divmod(minutes, 60)
        days, hours = divmod(hours, 24)
        tmp = (
            ((str(days) + "d, ") if days else "")
            + ((str(hours) + "h, ") if hours else "")
            + ((str(minutes) + "m, ") if minutes else "")
            + ((str(seconds) + "s, ") if seconds else "")
        )
        return tmp[:-2]

class HandleText:
    @staticmethod
    def input_str(union: Union[Message, InlineQuery]) -> str:
        """Get text before of a command..."""
        text = union.text if isinstance(union, Message) else union.query
        return " ".join(text.split()[1:])

class HandleBoolean:
    @staticmethod
    def get_ticket(state: bool):
        return "✅" if state is True else "☑️"