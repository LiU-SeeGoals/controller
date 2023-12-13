// import logo from './logo.svg';
// import './App.css';
// import React, { useState, useEffect, useRef } from "react";
// import Canvas from "./components/Canvas/Canvas.js";
// import Menu from "./components/Menu/Menu.js";

// function App() {
//   return (
//     <div className="App">
//       <Menu/>
//       <Canvas/>
//     </div>
//   );
// }

// export default App;

import logo from './logo.svg';
import './App.css';
import Canvas from "./components/Canvas/Canvas.js";
import Menu from "./components/Menu/Menu.js";
import React, { useState, useEffect, useRef } from "react";
import "./App.css";

function App() {
  const sidebarRef = useRef(null);
  const [isResizing, setIsResizing] = useState(false);
  const [sidebarWidth, setSidebarWidth] = useState(700);
  const initialGameState = [
    {"id": 0, "team":"blue", "x": 3300, "y": 5200, "speed_x":1, "speed_y":1, "selected": false},
    {"id": 1, "team":"blue", "x": 5, "y": 55, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 2, "team":"blue", "x": 5, "y": 65, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 3, "team":"blue", "x": 5, "y": 75, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 4, "team":"blue", "x": 5, "y": 85, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 5, "team":"blue", "x": 5, "y": 95, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 0, "team":"yellow", "x": 95, "y": 45, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 1, "team":"yellow", "x": 95, "y": 55, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 2, "team":"yellow", "x": 95, "y": 65, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 3, "team":"yellow", "x": 95, "y": 76, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 4, "team":"yellow", "x": 95, "y": 85, "speed_x":0, "speed_y":0, "selected": false},
    {"id": 5, "team":"yellow", "x": 95, "y": 95, "speed_x":0, "speed_y":0, "selected": false},

  ]
  const [gameState, setGameState] = useState(initialGameState);

  const startResizing = React.useCallback((mouseDownEvent) => {
    setIsResizing(true);
  }, []);

  const stopResizing = React.useCallback(() => {
    setIsResizing(false);
  }, []);

  const resize = React.useCallback(
    (mouseMoveEvent) => {
      if (isResizing) {
        setSidebarWidth(
          mouseMoveEvent.clientX -
            sidebarRef.current.getBoundingClientRect().left
        );
      }
    },
    [isResizing]
  );

useEffect(() => {
    const socket = new WebSocket('ws://localhost:8080/ws');

    socket.onmessage = (event) => {
      try {
        const gamestateFromServer = JSON.parse(event.data);
        if (gamestateFromServer.BlueTeam && gamestateFromServer.YellowTeam) {
            let gameStateCopy = [...gameState];

            for (let i = 0; i < gamestateFromServer.BlueTeam.length; i++) {
              gameStateCopy[i]["y"] = gamestateFromServer.BlueTeam[i]["PosX"] + 3300;
              gameStateCopy[i]["x"] = gamestateFromServer.BlueTeam[i]["PosY"] + 4800;
            }

            for (let i = 0; i < gamestateFromServer.YellowTeam.length; i++) {
              gameStateCopy[i + gamestateFromServer.BlueTeam.length].x = gamestateFromServer.YellowTeam[i].PosX;
              gameStateCopy[i + gamestateFromServer.YellowTeam.length].y = gamestateFromServer.YellowTeam[i].PosY;
            }

            setGameState(gameStateCopy);
            console.log(gameStateCopy);
        }
      } catch (e) {
        console.error('Error parsing message JSON', e);
      }
    };

    return () => {
      socket.close();
    };
  }, []); // Removed


  React.useEffect(() => {
    window.addEventListener("mousemove", resize);
    window.addEventListener("mouseup", stopResizing);
    return () => {
      window.removeEventListener("mousemove", resize);
      window.removeEventListener("mouseup", stopResizing);
    };
  }, [resize, stopResizing]);

  return (
    <div className="app-container">
      <div
        ref={sidebarRef}
        className="app-sidebar"
        style={{ width: sidebarWidth }}
        onMouseDown={(e) => e.preventDefault()}
      >
        <div className="app-sidebar-content">
        <Menu gameState={gameState}/>
        </div>
        <div className="app-sidebar-resizer" onMouseDown={startResizing} />
      </div>
      <div className="app-frame">
      <Canvas gameState={gameState} setGameState={setGameState}/>
      </div>
    </div>
  );
}

export default App;
