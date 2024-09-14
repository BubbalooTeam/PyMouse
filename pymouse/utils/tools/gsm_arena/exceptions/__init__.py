class GSMarenaManyRequests(Exception):
    """Raise if request is blocked due to 'Too Many Requests'."""

class GSMarenaBadRequest(Exception):
    """Raise if a 'BadRequest' error occurs."""

class GSMarenaCategoryError(Exception):
    """Raise if specified category is invalid."""

class GSMarenaDeviceNotFound(Exception):
    """Raise if the device requested by the user doesn't not exist or is unavailable on GSMarena."""

class GSMarenaPhoneInvalid(Exception):
    """Raise if 'phone' parameter is invalid."""

class GSMarenaFailIndex(Exception):
    """Raise if there is any failure to index the returned information."""