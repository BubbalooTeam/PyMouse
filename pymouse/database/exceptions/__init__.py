from json import JSONDecodeError

class LocalDatabaseJsonError(JSONDecodeError):
    """Raise if database format error."""

class LocalDatabaseNotFound(FileNotFoundError, FileExistsError):
    """Raise if database file doesn't exists"""