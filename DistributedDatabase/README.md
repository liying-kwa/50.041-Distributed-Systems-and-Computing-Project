# Distributed Database

Database to store enrolment information, written in golang.

## How to Setup
* To run the RingServer, run the command 'go run ringserver.go'
* To run a NodeServer, run the command 'go run nodeserver.go'

## Command Line Operations

### RingServer

| Command                   | Info                                                               |
|---------------------------|--------------------------------------------------------------------|
| help                      | Displays available commands.                                       |
| info                      | Displays information about the RingServer in PrettyPrint format.   |
| ring                      | Displays the ring structure in PrettyPrint format.                 |
| read <courseID>           | Carries out a read operation in the form of GET(courseID).         |
| write <courseID>, <count> | Carries out a write operation in the form of PUT(courseID, count). |

### NodeServer

| Command           | Info                                                                                                                                       |
|-------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| help              | Displays available commands.                                                                                                               |
| info              | Displays information about the NodeServer in PrettyPrint format.                                                                           |
| register <portNo> | Registers with the RingServer to add to the ring structure. If registration is successful, the NodeServer listens on the portNo specified. |


