# MyGoPractice

## How to setup the `golang v1.12` environment(based on ubuntu 19.10)

### 1. Install golang runtime
Install golang runtime v1.12 by running the following commands
```
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get install software-properties-common
sudo apt-get update 
sudo apt-get install golang-1.12
```
Now we need to add the go binary to our path. I like to put it in my `.profile` by typing:
```
vim ~/.profile
```
Then scrolling to the bottom and adding the following line
```
export PATH="$PATH:/usr/lib/go-1.12/bin"
```
Then use the source command to reread your profile.
```
source ~/.profile
```
Now type:
```
go version
```
And you should see that `Go 1.12` is installed.

### 2. Setup golang environment variables
Create a new folder under the user home, where to store all the golang related utility and source code in
```
mkdir -p $HOME/gopath
```
Add golang environment variables `GOROOT` and `GOPATH` into profile `.profile`
```
vim ~/.profile
```
Then scrolling to the bottom and adding the following line
```
export GOROOT="/usr/lib/go-1.12"
export GOPATH="$HOME/gopath"
export PATH="$PATH:$GOPATH/bin"
```
Then use the source command to reread your profile.
```
source ~/.profile
```
Now type:
```
go env
```
And you should see `GOROOT` and `GOPATH` have been set to the values above.

### 3. Install `libprotoc 3.9.1` for protobuffer foundation
Install the below libs
```
sudo apt-get install autoconf automake libtool curl make g++ unzip
```
Download `libprotoc 3.9.1` source code and install
```
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protobuf-all-3.9.1.zip
unzip protobuf-all-3.9.1.zip
cd protobuf-3.9.1
./configure
make
make check
sudo make install
sudo ldconfig
```
Now type:
```
protoc --version
```
And you should see that `libprotoc 3.9.1` is installed.

### 4. Install other golang utilities
Install golang protobuf v1.2.0
```
GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@v1.2.0
```
Install mockgen v1.4.3
```
GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.3
```
Install dlv v1.2.0
```
GO111MODULE=on go get github.com/go-delve/delve/cmd/dlv@v1.2.0
```
Install dep
```
go get -u github.com/golang/dep/cmd/dep
```

### 5. Install the 3rd party components
Install docker
```
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable edge"
sudo apt-get update
sudo apt-get install docker-ce
sudo usermod -aG docker ${USER}
```
**Note:** After doing the above, do remember to logout and then login again.

Install docker-compose
```
sudo apt-get install docker-compose
```

## Setup this project

### 1. Clone the code into the right place
Create a folder where to store the project code in
```
mkdir -p $GOPATH/src/github.com/yn295636
```
Clone the code
```
cd $GOPATH/src/github.com/yn295636
git clone https://github.com/yn295636/MyGoPractice
```

### 2. Install dependencies
Install the dependencies through dep
```
cd PROJ_FOLDER
dep ensure -v
```

### 3. Generate golang code for protobuf
Generate the golang code for defined protobuf `greeter_service.proto` and `sample_service.proto`
```
cd PROJ_FOLDER/proto/greeter_service
./build_proto.sh
cd PROJ_FOLDER/proto/sample_service
./build_proto.sh
```
Then, you have already been able to run the unit test under `app/sample_service`. You can try with below commands
```
cd PROJ_FOLDER/app/sample_service
go test -v
```

## Startup the whole project for production use

Startup the 3rd party components, include redis, mongo DB and mysql DB
```
cd PROJ_FOLDER
./launch_docker_for_develop.sh
```
Startup `greeter_service`
```
cd PROJ_FOLDER/app/greeter_service
go build
./greeter_service
```
Startup `sample_service`
```
cd PROJ_FOLDER/app/sample_service
go build
./sample_service
```
Startup `apigateway`
```
cd PROJ_FOLDER/app/apigateway
go build
./apigateway
```
Then try
```
curl "http://localhost:8081/multiply?a=2&b=3"
curl -X POST -H "Content-Type:application/json" -d '{"name":"Bobo"}' http://localhost:8081/greet
```
You should see corresponding logs from the services respectively.

## Run unit test for the services

### 1. Run unit test for `sample_service`
Run unit test for `sample_service`
```
cd PROJ_FOLDER/app/sample_service
go test -v
```

### 2. Run unit test for `greeter_service`
Run unit test for `greeter_service`
```
cd PROJ_FOLDER
./launch_docker_for_test.sh
cd PROJ_FOLDER/app/greeter_service
go test -v
cd PROJ_FOLDER
./unlaunch_docker_for_test.sh
```
Debug unit test for `greeter_service` if you want to see the intermediate result while running the test
```
cd PROJ_FOLDER
./launch_docker_for_test.sh
cd PROJ_FOLDER/app/greeter_service
dlv test -- -test.v
```
Then you can set breakpoint in `dlv` interactive command-line interface. For instance, you can set the breakpoint in the case regarding for operating mysql DB, and then check the data in mysql test DB (port: 13306, user: root, password: Mytest123!). For the detailed usage of `dlv`, please check the [doc](https://github.com/go-delve/delve/blob/master/Documentation/cli/README.md)

After the debug is finished, please teardown the test
```
cd PROJ_FOLDER
./unlaunch_docker_for_test.sh
```

### 3. Run unit test for `apigateway`
Run unit test for `apigateway`
```
cd PROJ_FOLDER/app/apigateway
go test -v ./test
```