# tctxto App

## How to kill process using a port

```
# pID is the process id
$ sudo lsof -i :2121
$ kill <pID> 
```

## How to retrieve submodules

```
# Initialize submodules
git submodule init
git submodule update --init --recursive

# Update submodules
git submodule update --remote
```

## How to prepare

```
# Make sure that you already install npm
make prepare
```

## How to build

```
make build-all
```

## How to run

```
# If you are in macOS
make run-macos

# If you are in linux
make run-linux

# If you are in windows
make run-windows
```

## How to run the released app

```
cd to/the/root/folder/containing/the/executables
./tctxtoapp # for linux and macOS
./tctxtoapp.exe # for windows

In the terminal, look for TCTXTO_PROXY_AO=http://192.168.1.10:2323

Then, enter the http url in your browser

Enjoy, have fun :D
```