from os import path, makedirs, remove
from shutil import rmtree

from pymouse import Config

class DownloadPaths:
    def __init__(self):
        self.filepath = Config.DOWNLOAD_PATH

    def Make_DownloadPath(self):
        if not path.exists(self.filepath):
            try:
                makedirs(
                    name=self.filepath,
                )
            except Exception:
                return "Fail in MAKE Directory.."
            return "Created directory with sucessfully.."
        return "Path already exists.."

    def Delete_DownloadPath(self):
        if path.exists(self.filepath):
            try:
                rmtree(
                    path=self.filepath,
                )
            except Exception:
                return "Fail in DELETE Directory.."
            return "Deleted directory and all content there.."
        return "Path doesn't exists.."

    def Delete_FileName(self, filename: str):
        if path.exists(filename):
            try:
                remove(filename)
            except Exception:
                return "Fail in DELETE FileName.."
            return "Deleted filename and all content there.."
        return "Path doesn't exists.."