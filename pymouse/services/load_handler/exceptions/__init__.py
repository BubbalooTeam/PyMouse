import re

class LoadModulesError(re.error, ImportError):
    """Raise if the modules are not loaded correctly"""