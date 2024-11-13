import torch
class Position:
    def __init__(self, x: float, y: float, z: float, angle: float) -> None:
        self.x = x
        self.y = y
        self.z = z
        self.angle = angle

    def __repr__(self) -> str:
        return f"x: {self.x:6.1f}, y: {self.y:6.1f}, z: {self.z:6.1f}, angle: {self.angle:5.2f}"
        
class Robot:
    def __init__(self, data: dict) -> None:
        self.id = data['Id']
        self.position = Position(data['X'], data['Y'], 0, data['Angle'])
        self.velocity = Position(data['VelX'], data['VelY'], 0, data['VelAngle'])

    def __repr__(self) -> str:
        return f"Robot{self.id}:\n      position: {self.position},\n      velocity: {self.velocity}"
    
    def to_torch(self):
        return torch.tensor([self.position.x, self.position.y, self.position.angle, self.velocity.x, self.velocity.y, self.velocity.angle]).float()
    
class Team:
    def __init__(self, data: list[dict], team_id: int) -> None:
        self.robots = {robot_data["Id"]: Robot(robot_data) for robot_data in data}
        self.team_id = team_id

    def __repr__(self) -> str:
        return f"Team{self.team_id}:\n    " + "\n    ".join([str(robot) for robot in self.robots.values()])
    
    def to_torch(self):
        team_positions = []
        for robot in self.robots.values():
            team_positions.append(robot.to_torch())
        team_positions = torch.stack(team_positions)
        return team_positions
    
class Ball:
    def __init__(self, data: dict) -> None:
        self.position = Position(data['PosX'], data['PosY'], data['PosZ'], 0)
        self.velocity = Position(data['VelX'], data['VelY'], data['VelZ'], 0)

    def __repr__(self) -> str:
        return f"Ball:\n    position: {self.position},\n    velocity: {self.velocity}"
    
    def to_torch(self):
        return torch.tensor([self.position.x, self.position.y, self.position.angle, self.velocity.x, self.velocity.y, self.velocity.angle]).float().unsqueeze(0)
           
class GameState:
    def __init__(self, data: dict) -> None:
        self.unknown_team = Team([robot for robot in data['RobotPositions'] if robot['Team'] == 0], 0)
        self.yellow_teams = Team([robot for robot in data['RobotPositions'] if robot['Team'] == 1], 1)
        self.blue_teams = Team([robot for robot in data['RobotPositions'] if robot['Team'] == 2], 2)
        self.ball = Ball(data['BallPosition'])

    def __repr__(self) -> str:
        return f"GameState: \n  {self.yellow_teams}, \n  {self.blue_teams}, \n  {self.ball}"
    


