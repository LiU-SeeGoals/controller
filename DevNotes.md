## Developer Diary: Defining Our Goals and Strategy
### Introduction
This entry serves as a guide to clarify and align the goals and expectations regarding the system being designed. In past discussions, several ideas were touched upon, but it's time to put them into writing for clarity and future reference. This will help onboard new team members more easily and ensure that the team stays aligned on the mission.

### Challenge the Status Quo
The current system is an idea and a draft of what we aim to build. However, it's important to challenge the status quo and continuously ask whether this is the best approach to achieve our goals. Openness to new ideas and flexibility in our methods are crucial. **Don’t be afraid to fail, be afraid not to try**. If something seems unclear or isn't working, it should be questioned—whether it's due to lack of clarity or an inherent flaw.

### Current System Overview
The current architecture provides a foundation, but the system lacks clearly defined goals. This document addresses that by setting objectives for how the system should evolve.

The primary objective is to create a decision engine capable of controlling all the players in a football game in real-time, with responsiveness being a top priority.

### The Dual-System Architecture
Two interconnected systems are needed:

* ***Rapid Response System***: A system that micromanages agents in the game, ensuring they follow rules and execute plans in real-time.
* ***High-Level Planning System***: A more time-flexible system that computes higher-level strategies and isn’t bound by the strict time constraints of the rapid response system.

The existing architecture effectively supports the rapid response part, but lacks support for strategy planning.

## Strategy
In terms of strategy, the team explored the possibility of hardcoding specific actions, such as:

**On offense:** positioning two agents higher up the field than the player with the ball.

**On defense:** blocking the line of sight for the opposing player with the ball.

However, manually designing all strategies would be time-consuming and not very flexible. To overcome this, a data-driven and adaptable system for planning is proposed.

### The Proposal: Data-Driven Strategy Planning
A height map-based data structure is suggested to assist in strategy generation. This system can dynamically represent the game by modeling opponents as "hills" and safe zones as "valleys." Here's how it would work:

Each opponent is at the top of a hill, and the further away from them, the lower the risk.
This representation is flexible enough to account for robot speed, mass, and acceleration.
By tracking robots' during a game a lot of these features can be discovered and thus the map can be dynamic and adaptable.
Using this data, the system can estimate where opponents and the ball are likely to be over a given time frame, introducing a layer of uncertainty. This flexible approach allows for modeling robots' speed, acceleration, and turning rates dynamically.


### Minimal Viable Product (MVP)

#### Offense
The first step is to create a system that can analyze the field and generate a map showing how long it would take each opponent to reach different points based on their current position and speed. Using this map, the system can then create a safe path toward the goal that minimizes the risk of the ball being intercepted by opponents. 

This system will be constantly updated in real-time as the game progresses to adapt to new situations.

**Requirements:**
* A system that generates a field map showing how quickly opponents can reach different areas.
* A system that calculates a safe path to the goal, avoiding interception by opponents.
* A system that moves the robots along the calculated path.
* A system that collects data on the robots and ball to keep the field map updated.
* A system that picks up the ball and positions it along the path.
* A system that passes the ball between robots.
* A system that catches/receives passes from teammates.
* A system that takes a shot on goal.

#### Defense
On defense, a similar approach can be used to protect against opponent attacks. The system should generate a map showing how quickly our robots can reach different areas on the defensive side of the field. The goal is to minimize the time it takes our robots to cover key areas, intercept the ball, and block the opponent's attack.

The system should also create paths for the robots to intercept the ball and block the movement of opponent robots, while preventing clear shots on our goal.

**Requirements:**
* A system that generates a field map for our defensive side, showing how quickly our robots can reach key areas.
* A system that moves the robots to minimize the time it takes to cover all important defensive areas.
* A system that calculates a path for our robots to intercept the ball.
* A system that moves the robots to block opponent robots and intercept passes.
* A system that blocks the opponent’s line of sight to our goal.


### Future Goals
The MVP will provide a solid foundation for the system, but there are several areas to improve upon in the future:

Currently, the height map representation serves as a useful way to model the game. It initially focuses on the time it takes for a robot to reach specific areas of the field, acting as a type of "risk planner." In the future, this approach can be expanded with more complex, data-driven features.

By starting with a height map that identifies risky regions, we can later integrate advanced techniques like deep learning or reinforcement learning to generate more sophisticated risk assessments. This height map allows the planner to execute quickly, enabling effective behavior from the start. Over time, the system can evolve by adding more complex features to create a more advanced height map.

To clarify, the computation of risky regions on the height map should be part of the high-level planner, which is not constrained by real-time demands. This planner can refine the height map over time. Meanwhile, the rapid response system should be able to operate with any height map and compute actions quickly, adapting when unexpected events occur, such as recovering a misdirected ball. The real-time system will handle the finer details of collision avoidance, passing, and ball recovery, while the high-level planner provides broader guidance on safe and risky areas.

In an ideal scenario, offense and defense would not be strictly separated. Instead, both could be represented within the height map, which the rapid response system could use to execute decisions in real-time. The system would adapt fluidly between offense and defense based on the situation.

This concept is still in its early stages, and more experimentation is needed to see how it performs when tested. Setting up the system to simulate games and letting the components interact will help refine these ideas further.
