# elnews backend server

### loads/my_config.go

##### 1. At server init, load _config.json_ files from config directory.

##### 2. It will dial and keep connection with DB at initial.


### api/posts.go

##### 1. Contains logic to fetch news posts data from db.


### main.go

##### 1. Main file loads configuration and routers, then start server
    
