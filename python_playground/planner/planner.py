from unified_planning.shortcuts import *

import networkx as nx
import random

from find_goal import get_model, predict

def generate_problem_coords(x_range,y_range, num_robots):
        
    jobs = []
    for robot in range(num_robots):
        origin_x = random.randrange(x_range)
        origin_y = random.randrange(y_range)
        
        destination_x = random.choice([x for x in range(x_range)
                                     if x != origin_x])
        destination_y = random.choice([x for x in range(y_range)
                                     if x != origin_y])                             
        jobs.append(((origin_x, origin_y), (destination_x, destination_y)))

    
    ball_location = (0, 0)
    
    return jobs, ball_location

def generate_problem_zones(num_zones, num_robots):
        
    jobs = []
    for robot in range(num_robots):
        origin = random.randrange(num_zones)

        destination = random.choice([x for x in range(num_zones)
                                     if x != origin])
                                 
        jobs.append((origin, destination))

    
    ball_location = random.randrange(num_zones)
    
    return jobs, ball_location


def get_planning_task_zones(jobs, ball_location, ball_goal, connected_zones):

    
    num_zones = len(list(connected_zones.keys()))
    
    Zone = UserType("Zone")
    Robot = UserType("Robot")

    zones = [up.model.Object(f"loc{i}", Zone)
                 for i in range(num_zones)]
    robots = [up.model.Object(f"robot{i}", Robot)
               for i in range(num_robots)]
    
  
    robot_at = up.model.Fluent("robot_at", BoolType(), r=Robot, l=Zone)
    connected = up.model.Fluent("connected", BoolType(), l_from=Zone, l_to=Zone)
    ball_at = up.model.Fluent("ball_at", BoolType(), l=Zone)
    has_ball = up.model.Fluent("has_ball", BoolType(), r=Robot)
    problem = up.model.Problem("PlaySoccer")

    problem.add_objects(zones)
    problem.add_objects(robots)
    problem.add_fluent(robot_at, default_initial_value=False)
    problem.add_fluent(connected, default_initial_value=False)
    problem.add_fluent(ball_at, default_initial_value=False)
    problem.add_fluent(has_ball, default_initial_value=False)

    problem.set_initial_value(ball_at(zones[ball_location]), True)
    problem.set_initial_value(has_ball(robots[0]), True)

    for robot_id, (origin, dest) in enumerate(jobs):
        r1 = robots[robot_id]
        problem.set_initial_value(robot_at(r1, zones[origin]), True)
        problem.add_goal(robot_at(r1, zones[dest]))

    for key in connected_zones.keys():
        
        for l in connected_zones[key]:
            
            problem.set_initial_value(connected(zones[key], zones[l]), True)

    # Define the move action
    move = up.model.InstantaneousAction("move", l_from=Zone, l_to=Zone, r=Robot)
    l_from = move.parameter("l_from")
    l_to = move.parameter("l_to")
    r = move.parameter("r")

    # Preconditions
    move.add_precondition(connected(l_from, l_to))  # Locations must be connected
    move.add_precondition(robot_at(r, l_from))  # Robot must be at the starting location
    move.add_precondition(Not(Equals(l_to, 9))) # Cannot move inside goal 
    move.add_precondition(Not(has_ball(r)))

    # Effects
    move.add_effect(robot_at(r, l_to), True)  # Robot arrives at the target location
    move.add_effect(robot_at(r, l_from), False)  # Robot leaves the starting location

       
    # Define the action
    pass_ball = up.model.InstantaneousAction("pass_ball", bl_from=Zone, bl_to=Zone, r_from=Robot, r_to=Robot)
    bl_from = pass_ball.parameter("bl_from")
    bl_to = pass_ball.parameter("bl_to")
    r_from = pass_ball.parameter("r_from")
    r_to = pass_ball.parameter("r_to")

    #Preconditions
    pass_ball.add_precondition(ball_at(bl_from))  # Ball must be at the starting location
    pass_ball.add_precondition(robot_at(r_from, bl_from))
    pass_ball.add_precondition(robot_at(r_to, bl_to))
    pass_ball.add_precondition(Not(Equals(r_to, r_from)))
    pass_ball.add_precondition(Not(Equals(bl_to, zones[9]))) # Don't pass into goal

    # Effects
    pass_ball.add_effect(ball_at(bl_from), False)  # Ball is no longer at the starting location
    pass_ball.add_effect(ball_at(bl_to), True)  # Ball is now at the target location

    # Define the shoot action 
    shoot = up.model.InstantaneousAction("shoot", r=Robot)
    r = shoot.parameter("r")

    #Preconditions
    shoot.add_precondition(has_ball(r))  # Ball must be at the starting location
    shoot.add_precondition(robot_at(r, zones[7])) # Shoot right in front of goal

    # Effects
    shoot.add_effect(has_ball(r), False)  # Ball is no longer at the starting location
    shoot.add_effect(ball_at(zones[9]), True)  # Ball is now at the target location

    # Define the shoot action 
    dribble = up.model.InstantaneousAction("dribble", r=Robot, rl_from=Zone, rl_to=Zone)
    r = dribble.parameter("r")
    rl_from = dribble.parameter("rl_from")
    rl_to = dribble.parameter("rl_to")

    # Preconditions
    dribble.add_precondition(has_ball(r)) # Ball must be at the starting location
    dribble.add_precondition(robot_at(r, rl_from))
    dribble.add_precondition(connected(rl_to, rl_from))

    # Effects
    dribble.add_effect(ball_at(rl_to), True)  
    dribble.add_effect(ball_at(rl_from), False)
    dribble.add_effect(robot_at(r, rl_from), False)  
    dribble.add_effect(robot_at(r, rl_to), True) 

    # Add action to the problem
    problem.add_action(dribble)
    problem.add_action(shoot)
    problem.add_action(pass_ball)
    problem.add_action(move)
    
    problem.add_goal(ball_at(zones[ball_goal]))  
    #print(problem.actions)
   
    return problem 


if __name__ == "__main__":

    num_zones = 9
    num_robots = 6 
    
    field = {0:[1,4,3], 1:[0,2,3,4,5], 2:[1,4,5], 3:[0,1,4,6,7,9], 4:[i for i in range(0,9)], 5:[1,2,4,7,8], 
    6:[3,4,7], 7:[3,4,5,6,8], 8:[4,5,7], 9:[7]}

    jobs, ball_location = generate_problem_zones(num_zones, num_robots)

    #model = get_model("model.pt")
    #ball_goal =  predict(model, jobs).item()
    ball_goal = 9

    print("Jobs:")
    for robot_id, (origin, destination) in enumerate(jobs):
        print(f"Robot {robot_id} from location {origin} to location {destination}")
    print(f"The ball is at location {ball_location} and should go to {ball_goal}")


  
    problem = get_planning_task_zones(jobs, ball_location, ball_goal, field)
    with OneshotPlanner(name="fast-downward") as planner:        
        result = planner.solve(problem)
        if result.status == up.engines.PlanGenerationResultStatus.SOLVED_SATISFICING:
            print(f"Fast Downward returned {result.plan}")
        else:
            print("No plan found.")

