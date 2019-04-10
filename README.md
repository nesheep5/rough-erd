# rough-erd
This tool creates a rough ER diagram.

## Install
```
go get -u github.com/nesheep5/rough-erd
```

# Usage
```
NAME:
   rough-erd - This tool creates a rough ER diagram.

USAGE:
   rough_erd [global options] command [command options] [arguments...]

VERSION:
   v1.0.0

COMMANDS:
     make, m  make ER diagram.
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

# Example
```
> cd $[GOPATH]/src/github.com/nesheep5/rough-erd
> docker-compose up -d --build
> rough_erd  make -P 23306 -u root -p rough-erd -protocol tcp -n test
create MySQL DB...
connection: root:rough-erd@tcp(127.0.0.1:23306)/test?maxAllowedPacket=0
created MySQL DB
-----------------------------

@startuml
title Sample ERDiagram
entity "companies" {
id
}
entity "offices" {
id
company_id
}
offices -- companies :company_id
@enduml

-----------------------------
Open → http://www.plantuml.com/plantuml/uml/UDhYSYWkIImgAStDuIh9BCb9LGXEp2t8ILLm3NB9J4mlIipbIiqhoIofL51ApiyjICpBJ2rMKgZcoapXgeNBvAUbPIR3nI7gAkF1Ig1I2hgw2d3z2bP8IXnIyr90bWC2003__x3bBqm0

```
Open: http://www.plantuml.com/plantuml/uml/UDhYSYWkIImgAStDuIh9BCb9LGXEp2t8ILLm3NB9J4mlIipbIiqhoIofL51ApiyjICpBJ2rMKgZcoapXgeNBvAUbPIR3nI7gAkF1Ig1I2hgw2d3z2bP8IXnIyr90bWC2003__x3bBqm0
