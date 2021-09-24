# go-learning
programs to learn go

goroutines and channels usage and test

A simple client and server program to test net package. It will act like ping - pong and continue in a loop.

Start the server program first, it will be listening on the port configured.

Then you may start one or more client programs to test.


example:-
$go run server.go
Starting the server ...



$go run client.go 
This is Dial ID: 0
Received-> 0
Sending to the server = 0
This is Dial ID: 1
This is Dial ID: 2
Received from server =  1
Received-> 1
Sending to the server = 1
This is Dial ID: 3
This is Dial ID: 4
This is Dial ID: 5
Received from server =  2
This is Dial ID: 6
Received-> 2
Sending to the server = 2
This is Dial ID: 7
Received from server =  3
This is Dial ID: 8
Received-> 3



