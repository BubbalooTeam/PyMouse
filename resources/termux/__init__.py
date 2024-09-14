#    PyMouse (Telegram BOT Project)
#    Copyright (c) 2022-2024 - BubbalooTeam

#    This program is free software: you can redistribute it and/or modify
#    it under the terms of the GNU Affero General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.

#    You should have received a copy of the GNU Affero General Public License
#    along with this program.  If not, see <https://www.gnu.org/licenses/>.


import os

class TermuxUtils:
    PyMouseLogo = """
     _____       __  __
    |  __ \     |  \/  |
    | |__) |   _| \  / | ___  _   _ ___  ___
    |  ___/ | | | |\/| |/ _ \| | | / __|/ _ \\
    | |   | |_| | |  | | (_) | |_| \__ \  __/
    |_|    \__, |_|  |_|\___/ \__,_|___/\___|
            __/ |
           |___/
                     """

    TermuxInstallerMessage = """
    • PyMouse Installer via Termux:

        • Initializing the Installation (Do Not Turn Off Your Device During Installation...).

        • It is possible that the process will take a while.
                     """

    PIP_COMMAND = "pip install -Ue ."

    APT_PACKAGES = [
        "neofetch",
        "ffmpeg",
    ]

class Termux:
    def print_empty(self):
        print()

    def install_apt(self, packages: list):
        self.print_empty()
        print("[...] - Installing packages using APT...\n")
        apt_packages_name = " ".join(TermuxUtils.APT_PACKAGES)
        os.system("apt install {packages} -y".format(packages=apt_packages_name))
        self.print_empty()
        print("[✔] - Installed packages with APT...\n")

    def install_pip(self, pip_cmd: str):
        self.print_empty()
        print("[...] - Installing pip requirements...\n")
        os.system(pip_cmd)
        self.print_empty()
        print("[✔] - Installed packages with PIP/Python3...\n")

    def run_bot(self):
        self.print_empty()
        print("[✔] - Done! Dependencies of bot is now installed, Starting PyMouse...\n")
        os.system("python -m pymouse")

    def installer_bot(self):
        print(TermuxUtils.PyMouseLogo)
        self.print_empty()
        print(TermuxUtils.TermuxInstallerMessage)
        # Initialize instalation
        self.install_apt(TermuxUtils.APT_PACKAGES)
        self.install_pip(TermuxUtils.PIP_COMMAND)
        self.run_bot()

Termux().installer_bot()