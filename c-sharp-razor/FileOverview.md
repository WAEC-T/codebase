# highlevel overview of the initial code for the project

### control.sh

The script that is called when alteration of the server-state is wanted.

### flag_tool

The script that inspects and interacts with the database *(minitwit.db)* when *control.sh* is given certain flags, through the command line.

### flag_tool.c

The source code that compiled makes the *flag_tool* executable.

### Makefile *(Not currently in use)*

* Init the database
* Build/Compile the *flag_tool*
* Remove the compiled *flag_tool*

### minitwit_tests.py

Test suit for the program.

### minitwit.db

The database file.

### minitwit.py

Creates the entire web app and interacts with the database.

### minitwit.pyc

Compiled bytecode file that are generated by the Python interpreter when *minitwit.py* script is executed.

### schema.sql

The sql file containing the sql code for the db tables.