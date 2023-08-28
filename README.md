# JPROJ

## 1. What's JPROJ?
JPROJ is a simple java build system.

## 2. Usage

### 2.1 Create a new project
```sh
$ jproj createproject -n "newproject"
```
This command will create `newproject`, `newproject/src` and the `newproject/target` directories. Besides creating those
directories, the `createproject` action will also create the file `newproject/jproj.json`, which contains information about
the project (*such as the project's name, and the project's target directory*).

### 2.2 Build a project
```sh
$ jproj build
```
Compiles the java source files and outputs it to the target's directory (*which is specified in the `newproject/jproj.json` file.*)

### 2.3 Add new elements to the project (*classes, interfaces, enums*)
```sh
$ jproj new -n com.company.app.App -t class
```
This will create the file `newproject/src/me/company/app/App.java`. This file already contains an empty Java class.
```java
package com.company.app;

public class App {
}
```

### 2.4 Clean the compiled classes in the target directory
```sh
$ jproj clean
```

## 3. How to compile JPROJ from sources
```sh
$ make build
```
The binary is located at `./bin/jproj`