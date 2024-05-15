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

    - Termux installer. . . Initializing instalation!
                     """
    APT_PACKAGES = [
        "neofetch",
        "ffmpeg",
    ]

print(TermuxUtils.PyMouseLogo)
print()
print("[...] - Installing packages using APT...")
apt_packages_name = " ".join(TermuxUtils.APT_PACKAGES)
os.system("apt install {packages} -y".format(apt_packages_name))
print()
print("[✔] - Installed packages with APT...")
print()
print("[...] - Installing pip requirements...")
os.system("pip install -Ue .")
print()
print("[✔] - Installed packages with PIP/Python3...")
print()
print("[✔] - Done! Requirements is now installed, Starting PyMouse...")
os.system("python3 -m pymouse")