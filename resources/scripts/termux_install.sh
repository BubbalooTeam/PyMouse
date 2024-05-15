printf "Updating system to install BOT, Please wait. . .\n\n"
pkg update -y
apt update
apt update -y

python_not_installed="$(python -c 'exit()')"
# Check if python is installed
if [python_not_installed]
then
    printf "Python3 not Found, Installing Python. . .\n"
    pkg install python3 python3-pip -y
fi

printf "Ok!! Instalation is terminated! Go to installation. . .\n\n\n"
python3 resources/scripts/termux/__init__.py