from flask import Flask, request
from pprint import pprint
import json

app = Flask(__name__)

@app.route("/")
def hello():
    print("a hello was received")
    return r"Hello back"

@app.post("/update_game_state/")
def update_game_state():
    print("a post was received")
    game_state = request.get_json()
    # = json.loads(raw_json)
    print(type(game_state))
    pprint(game_state)
    return r"Plan received"

@app.get("/plan/")
def get_plan():
    return r"{}"