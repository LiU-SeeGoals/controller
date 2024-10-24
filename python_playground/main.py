from flask import Flask, request, jsonify
from flask_cors import CORS
from gamestate import GameState
from models import SeeGoalsDNN
import torch
import pathlib
import time

app = Flask(__name__)
CORS(app)

model = SeeGoalsDNN(num_players_per_team=2)
if pathlib.Path(model.path).exists():
    print("Loading model")
    model.load_state_dict(torch.load(model.path))
    model.eval()
calls = 0
start = time.perf_counter()
@app.route('/slowBrainBlue', methods=['POST'])
def slow_brin_blue():
    gamestate = GameState(request.get_json())
    my_team = gamestate.blue_teams
    enemy_team = gamestate.yellow_teams
    ball = gamestate.ball
    scale = 100
    output = model(my_team.to_torch().unsqueeze(0), enemy_team.to_torch().unsqueeze(0), ball.to_torch().unsqueeze(0)) * scale
    instructions = []

    for pred, robot in zip(output[0], my_team.robots.values()):
        dx, dy = pred.tolist()
        dest = [robot.position.x + dx, robot.position.y + dy, 0, 0]

        instructions.append({
            "Id": robot.id,
            "Position": dest
        })
    

    plan = {
        "Instructions": instructions,
    }
    global calls
    global start

    end = time.perf_counter()
    calls += 1
    cps = calls / (end - start)
    print(f"CPS: {cps:.2f}")
    return jsonify(plan), 200

@app.route('/slowBrainYellow', methods=['POST'])
def slow_brin_yellow():
    gamestate = GameState(request.get_json())
    my_team = gamestate.yellow_teams
    enemy_team = gamestate.blue_teams
    ball = gamestate.ball
    scale = 100
    output = model(my_team.to_torch().unsqueeze(0), enemy_team.to_torch().unsqueeze(0), ball.to_torch().unsqueeze(0)) * scale
    instructions = []

    for pred, robot in zip(output[0], my_team.robots.values()):
        dx, dy = pred.tolist()
        dest = [robot.position.x + dx, robot.position.y + dy, 0, 0]

        instructions.append({
            "Id": robot.id,
            "Position": dest
        })
    

    plan = {
        "Instructions": instructions,
    }
    print(plan)
    return jsonify(plan), 200

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000, debug=True)