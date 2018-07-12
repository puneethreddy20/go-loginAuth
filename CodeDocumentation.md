#Code Documentation

The backend is written  in Golang. It has the following functionalities: 
1. Users can register, login with a username and password, and view a custom greeting page.
2. When a user visits creating page, the backend calls the GreetingAPI with a signed token Json signed web tokens which is valid for 30 mins.
3. Once the login, it creates a cookie for further authentication and is valid for 6 hours.
4. It uses sqlite3 database for storing users information such as useremail, username, and password stored as hash (generated using bcrypt library).


I haven't choosen any big frameworks such as gorilla, martini, because net/http package provides more than enough functions. And I didn't want to make the application heavy by importing many packages.

The website checks for CSRF attacks, XSS and uses secure cookies for authentication purposes, prevents sql injection by preparing statements before hand.
Using mutex/locks to avoid access at same time.

JWT token is signed with a shared secret key, There are basically two ways for signing Token.

One is using shared secret key. This key is used to generate signed token and also validate at other end.

Other is signing with private certificate and validating with public certificate.

Its better to use shared key, I have assumed the key is already present with two services. 

As for unit tests, i have implemented test values for verifying the handlers. It is possible mock the database functions for testing, but as sqlite3 is lightweight, creating dummy test db should be fine.

config.yml consists database url, GreetingAPI url, templates path. It would be easy just to change config file for frequently used values.



#Futher work:
1. With certificates, TLS/SSL (https) can be implemented.
2. JWT with public and private certificates.
3. Add Forgot password, Forgot username.
4. Implement Email Functions.