![Team logo](./images/logo.png)

# See Goals Robot Controller

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
* `GRSIM_COMMAND_LISTEN_PORT` - grsim command listen port
* `GC_PUBLISH_ADDR` - multicast IP used by game controller
* `GC_PUBLISH_PORT` - publish port used by game controller
* `WEB_VISION_UI_PORT` - port on host machine for SSL vision UI when running docker
* `WEB_GC_UI_PORT` - port on host machine for game controller UI when running docker

## Docker environment
The docker environment should be used for local development. It uses grsim to simulate the game.

To start the environment:
```sh
./scripts/compose_up.sh
```

This will start the docker environment (in detached mode). The Seegoals controller is meade to be run from inside the container. The controller container can be entered by:
```sh
./scripts/enter.sh
```

Taking down the environment is done with
```sh
./scripts/compose_down.sh
```

Note, since the Seegoals controller and grsim run in different containers, grsim is **not** reachable at `127.0.0.1`. Instead, you should set `GRSIM_ADDR = "grsim"` in `.env`.

(Docker creates aliases for containers in the same network, so `"grsim"` will resolve to the grsim container IP.)

## Simulation
Seegoals uses **grsim** for game simulation. It's possible to run grsim in headless mode (no UI) by using the docker environment.

### Actions
Send actions to GrSim by using AddActions() owned by GrsimClient. AddActions takes a slice of actions, ordered by robot id, and translates the action to parameters accepted by GrSim.

Supported actions:
- Dribble (cant be used simultaneously as other actions)
- Kick
- Move
- Stop (useless)
