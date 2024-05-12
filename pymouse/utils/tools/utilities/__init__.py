from typing import Callable
from functools import partial, wraps
from asyncio import get_event_loop

class UtilsDecorator:
    def aiowrap(self):
        def decorator(func: Callable) -> Callable:
            @wraps(func)
            async def wrapper(*args, loop=None, executor=None, **kwargs):
                if loop is None:
                    loop = get_event_loop()
                pfunc = partial(func, *args, **kwargs)
                return await loop.run_in_executor(executor, pfunc)
            return wrapper
        return decorator