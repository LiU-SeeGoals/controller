![Team logo](./images/logo.png)

# See Goals Robot Controller

## Project Structure

Project structure is based on: https://github.com/golang-standards/project-layout
If you ever wonder where to put new files, please refer to it.

## Code standard
This project uses the [SeeGoals Go standard](https://github.com/LiU-SeeGoals/wiki/wiki/1.-Processes-&-Standards#seegoal-%F0%93%85%B0---go-coding-standard).

## Compiling/Building

The project can be compiled by running `/scripts/build.sh`. This generates the executable in `/build` folder.

## GrSim
### Actions
Send actions to GrSim by using AddActions() owned by GrsimClient. AddActions takes a slice of actions, ordered by robot id, and translates the action to parameters accepted by GrSim.

Supported actions:
- Dribble (cant be used simultaneously as other actions)
- Kick
- Move
- Stop (useless)
