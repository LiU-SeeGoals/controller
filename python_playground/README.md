# Python SlowBrain

## running the server 


#### venv

```bash
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

#### run the server

```bash
python3 main.py
```

## Communication 

#### Input 

a post request is sent to the python server from the go server with the json gamestate


#### Output

a post response is sent to the go server with the json gameplan


