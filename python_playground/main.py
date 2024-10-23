from flask import Flask, request, jsonify
from flask_cors import CORS
from gamestate import GameState

app = Flask(__name__)
CORS(app)

@app.route('/slowBrain', methods=['POST'])
def slow_brin():
    gamestate = GameState(request.get_json())
    print(gamestate)
    plan = {
        "Instructions": [
            {"Id": 0, 
             "Position": [0, 0, 0, 0],
            },
            {"Id": 1,
             "Position": [0, 0, 0, 0],
            },
        ]

    }
    return jsonify(plan), 200

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000, debug=True)