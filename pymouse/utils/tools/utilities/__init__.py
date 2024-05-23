import time

from hydrogram.types import Message

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
    def input_str(m: Message) -> str:
        """Get text before of a command..."""
        return " ".join(m.text.split()[1:])

class HandleBoolean:
    @staticmethod
    def get_ticket(state: bool):
        return "✅" if state is True else "☑️"