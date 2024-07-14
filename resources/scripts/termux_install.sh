echo "[...] - Updating system to install BOT, Please wait. . .\n\n"
pkg update -y
apt update
apt update -y

# Check if python is installed
if ! command -v python &> /dev/null
then
    echo "[...] - Python not Found, Installing Python. . .\n"
    pkg install python python-pip -y
    echo "[✔] - Python and tools installed with successfully!\n"
fi

echo "[✔] Updates of Termux is terminated! Initializing Install PyMouse Tool. . .\n\n\n"
python resources/termux/__init__.py