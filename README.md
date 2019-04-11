# rough-erd
This tool creates a rough ER Diagram.  
It uses the ID to create a ERD. (Not using FOREIGN KEY)  
It can make ERD Text([PlantUml](http://plantuml.com) Format), PNG, SVG, and [online editor](http://www.plantuml.com/plantuml/) URL.

## Sample ERD
![uml.png](https://github.com/nesheep5/rough-erd/blob/master/example/uml.png)

## Relation Role
This tool uses the ID to relate.  Not using FOREIGN KEY.
Example, `salaries` table has `employee_id`, `salaries` related `employees`.   
(example sql:  `SELECT * FROM employees JOIN salaries ON employees.id = salaries.employee_id;`)

## Install
```bash
go get -u github.com/nesheep5/rough-erd/cmd/rough-erd
```

## Usage
```bash
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
### make command
```bash
NAME:
   rough_erd make - make ER diagram.

USAGE:
   rough_erd make [command options] [arguments...]

OPTIONS:
   --dbtype value              database type (default: "mysql")
   --user value, -u value      database user
   --password value, -p value  database password
   --host value, -H value      database host (default: "127.0.0.1")
   --port value, -P value      database port (default: 3306)
   --protocol value            database protocol
   --name value, -n value      database name
   --output value, -o value    output style [text, url, png, svg]  (default: "text")
```

## Example
### Run TestDB
```bash
> cd $[GOPATH]/src/github.com/nesheep5/rough-erd/example
> docker-compose up -d --build 
> mysql -P23306 -uroot -prough-erd --protocol=TCP -Dtest
mysql: [Warning] Using a password on the command line interface can be insecure.
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 8
Server version: 5.7.25 MySQL Community Server (GPL)

Copyright (c) 2000, 2019, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show tables;
+----------------+
| Tables_in_test |
+----------------+
| departments    |
| dept_emp       |
| dept_manager   |
| employees      |
| salaries       |
| titles         |
+----------------+
6 rows in set (0.00 sec)
```

### make UML Text
```bash
> rough_erd make -P 23306 -u root -p rough-erd -protocol tcp -n test -o text
-----------------------------

@startuml
title ER Diagram
entity "departments" {
id
}
entity "dept_emp" {
id
employee_id
department_id
}
entity "dept_manager" {
id
employee_id
department_id
}
entity "employees" {
id
}
entity "salaries" {
id
employee_id
}
entity "titles" {
id
employee_id
}
dept_emp -- employees :employee_id
dept_emp -- departments :department_id
dept_manager -- employees :employee_id
dept_manager -- departments :department_id
salaries -- employees :employee_id
titles -- employees :employee_id
@enduml

-----------------------------
```

### make UML URL
```bash
> rough_erd make -P 23306 -u root -p rough-erd -protocol tcp -n test -o url
Open → http://www.plantuml.com/plantuml/uml/UDgKaKsgmp0CXFSwXSW-5yWgY_Skq0i4WKKGM6wmrKKelNlfZnstAIvTppVZ6Gl6P1Jjf1vCp3F-7_1FQ8wamC74LkmSBnHDELZgy0pYu59hDh4kJu5rySULUH87cstQMvG2pHn_i6Lcto6HfoX5gCCswBxkCv8tODzZUGM7jr85gRu3XzUszRHlQHNMICpR6ccFPGrWvE1k1xu6003__qByi0K0
```
Open → http://www.plantuml.com/plantuml/uml/UDgKaKsgmp0CXFSwXSW-5yWgY_Skq0i4WKKGM6wmrKKelNlfZnstAIvTppVZ6Gl6P1Jjf1vCp3F-7_1FQ8wamC74LkmSBnHDELZgy0pYu59hDh4kJu5rySULUH87cstQMvG2pHn_i6Lcto6HfoX5gCCswBxkCv8tODzZUGM7jr85gRu3XzUszRHlQHNMICpR6ccFPGrWvE1k1xu6003__qByi0K0

### make UML png
```bash
> rough_erd make -P 23306 -u root -p rough-erd -protocol tcp -n test -o png > uml.png
> open uml.png
```
![uml.png](https://github.com/nesheep5/rough-erd/blob/master/example/uml.png)

### make UML svg
```
> rough_erd make -P 23306 -u root -p rough-erd -protocol tcp -n test -o svg > uml.svg
> open uml.svg -a /Applications/Google\ Chrome.app/
```
![uml.svg](https://github.com/nesheep5/rough-erd/blob/master/example/uml.svg)
