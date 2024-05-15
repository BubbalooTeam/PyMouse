printf "Updating system to install BOT, Please wait. . .\n\n"
pkg update -y
apt update
apt update -y

# Check if python is installed
if ! command -v python &> /dev/null
then
    printf "Python not Found, Installing Python. . .\n"
    pkg install python python-pip -y
fi

printf "Ok!! Updates of Termux is terminated! Proceed to the installation of PyMouse. . .\n\n\n"
python resources/termux/__init__.py