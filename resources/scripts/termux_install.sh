printf "Updating system to install BOT, Please wait. . .\n\n"
pkg update -y
apt update
apt update -y

python_not_installed="$(python3 -c 'exit()')"
# Check if python is installed
if [ python_not_installed ]
then
    printf "Python3 not Found, Installing Python. . .\n"
    apt install python3 python3-pip -y
fi

printf "Ok!! Updates of Termux is terminated! Go to installation of PyMouse. . .\n\n\n"
python3 resources/termux/__init__.py