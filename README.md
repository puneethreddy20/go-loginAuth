# How to build and run the application:

To run the application you would need docker installed (or) can use build and run directly. Install docker first.
Unzip the package

1. Using docker:
```
$ cd go-loginAuth

//build a image with dockerfile in the directory which is named as backend
$ sudo docker build -t backend .

//Run the container. Which starts the backend service on port 8080
$ sudo docker run -d -p 8080:8080 backend

```

2. Build and run:
```
$ cd go-loginAuth

//Run go build
$ go build

//Run the executable file generated by go build 
$ sudo ./go-loginAuth

```

On the host machine open browser and visit http://localhost:8080 .

It displays a page with login form and Register form. You can register for new user and then login.

Default Credentials:
```
username: test

password: test

```
If those doesn't work, please create/register.

The mechanism is explained in code documentation.

