# SCTUN
## multiplexing socks5 back-connect proxy

**how it works, short diagram**
![screenshot](https://sun9-3.userapi.com/c813024/v813024343/93922/wqX4dZxtF1Y.jpg)


build use Makefile
````bash
make linux
````
build to ./build/osx/

start
````bash
#start master stun server
./stun --addr 0.0.0.0:8080

#start slave ctun client
./ctun --master 0.0.0.0:8080
````

how to use. 
After success connect client to server, stun server print stdout port
````bash 
....
starting new socks5 proxy-server on 0.0.0.0:10001
````
you need use in socks5 proxy client:
ip - as stun server addr 
port - as stdout port (example :10001) 

Example curl
````bash
curl --socks5 127.0.0.1:10001 check-host.net/ip
````

