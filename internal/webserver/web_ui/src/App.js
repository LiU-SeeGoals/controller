import React, { useState, useEffect } from 'react';

function App() {
    const [robotPosition, setRobotPosition] = useState(null);

    useEffect(() => {
        const socket = new WebSocket('ws://localhost:8080/ws'); // Adjust URL as neede
        socket.onmessage = (event) => {
            const gameState = JSON.parse(event.data);
            console.log(gameState)

            // Assuming gameState has a structure where the first robot's position
            // can be accessed like this. Adjust according to your actual structure.
            if (gameState && gameState.blue_team && gameState.blue_team.length > 0) {
                const firstRobot = gameState.blue_team[0];
                setRobotPosition(firstRobot.pos); // Assuming 'pos' is the position object
            }
        };

        return () => {
        };
    }, []);

    return (
      <>
        <div>
            {robotPosition && <p>First Robot Position: X: {robotPosition.x}, Y: {robotPosition.y}</p>}
        </div>
        <div>
          hello
        </div>
      </>
    );
}

export default App;