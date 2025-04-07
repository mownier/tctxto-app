# tctxto App

```
# Getting the server, client, and proxy codebase
$ make init

# Start server, client, and proxy servers
# Client might have a delay in getting dependencies
# Re-run this if client is not starting
$ make run-all
```

## How to kill process using a port

```
# pID is the process id
$ sudo lsof -i :2121
$ kill <pID> 
```
