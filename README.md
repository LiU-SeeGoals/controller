![Team logo](./images/logo.png)

# See Goals Robot Controller

## Overview
### What is the purpose of this repo?
This is the repo containing the AI of the project. All the stratergies on how the robots will move and respond to opponents actions. This repo also is reponsible for the connections of our internal repos and external SSL repos.  

## Setup 🚀

### starting docker
The prefered developmen environment is docker. To build the project run following command in project root.
```
docker compose build
```
```
docker compose up
```

### webbpage
The webpage can be found by typing following addres in your webbrowser after starting the containers:
```
http://localhost:5173/
```

### running main script
Now the containers should be running. To run a program you need to enter the controller container. This can be done by running:
```
docker ps
```
Take note of the "container ID"
Then run this command to enter the container:
```
docker exec -it {first 3 letters of container ID} sh
``` 
Now you are inside the container. To start the controller main program, go to 
```
~/cmd
``` 
And then run:
```
go run main.go
``` 

Now the main program have been started. 


## Project Structure

Project structure is based on: https://github.com/golang-standards/project-layout
If you ever wonder where to put new files, please refer to it.

## Code standard
This project uses the [SeeGoals Go standard](https://github.com/LiU-SeeGoals/wiki/wiki/1.-Processes-&-Standards#seegoal-%F0%93%85%B0---go-coding-standard).

## Compiling/Building

The project can be compiled by running `/scripts/build.sh`. This generates the executable in `/build` folder.

## Environment
Environment configuration cen be found in `.env`. This file is automatically loaded by the `config` package by the controller. It's strongly advised to use this file instead of hard coded solutions. Apart from the controller, the docker environment loads the `.env` file into its containers.

Following are the most important environment variables:

* `ENVIRONMENT` - environment flag to indicate what setup is being used
* `SSL_VISION_MULTICAST_ADDR` - multicast IP used by SSL vision
* `SSL_VISION_MAIN_PORT` - port used for tracking, detection, and geometry packets
* `GRSIM_ADDR` - grsim IP address
* `SIM_COMMAND_LISTEN_PORT` - sim command listen port
* `GC_PUBLISH_ADDR` - multicast IP used by game controller
* `GC_PUBLISH_PORT` - publish port used by game controller
* `WEB_VISION_UI_PORT` - port on host machine for SSL vision UI when running docker
* `WEB_GC_UI_PORT` - port on host machine for game controller UI when running docker

<!-- ## Docker environment
The docker environment should be used for local development. It uses grsim to simulate the game.

To start the environment:
```sh
./scripts/compose_up.sh
```

This will start the docker environment (in detached mode). The Seegoals controller is meant to be run from inside the container. The controller container can be entered by:
```sh
./scripts/enter.sh
```

Taking down the environment is done with
```sh
./scripts/compose_down.sh
``` -->