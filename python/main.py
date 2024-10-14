from flask import Flask, request
import threading
from pprint import pprint
import json
from time import sleep
import logging
from collections import namedtuple

RobotPosition = namedtuple('RobotPosition', 'id team x y vel_x vel_y')
BallPosition = namedtuple('BallPosition', 'x y vel_x vel_y')

global game_state
game_state = None

app = Flask(__name__)

@app.post("/update_game_state/")
def update_game_state():
    global game_state
    game_state = request.get_json()
    # pprint(game_state)
    return r"Plan received!"

log = logging.getLogger('werkzeug')
log.disabled = True

flask_thread = threading.Thread(target=lambda: app.run(host='0.0.0.0', port='5000', ))
flask_thread.start()

def parse_game_state():
    robot_positions = []
    for robot_d in game_state['RobotPositions']:
        vs = []
        # TODO: Parse angles
        for key in ['Id', 'Team', 'PosX', 'PosY', 'VelX', 'VelY']:
            vs.append(robot_d[key])

        robot_position = RobotPosition(*vs)
        robot_positions.append(robot_position)

    ball_d = game_state['BallPosition']
    vs = []
    # TODO: Parse angles? Maybe? The ball is approximately round?
    for key in  ['PosX', 'PosY', 'VelX', 'VelY']:
        vs.append(ball_d[key])
    ball_position = BallPosition(*vs)

    return robot_positions, ball_position

while True:
    print('tick')
    if game_state is not None:
        robot_positions, ball_position = parse_game_state()
        print(robot_positions)
        print(ball_position)
    print('end of tick')
    sleep(1)