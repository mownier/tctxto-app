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