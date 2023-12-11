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
        <Menu/>
        </div>
        <div className="app-sidebar-resizer" onMouseDown={startResizing} />
      </div>
      <div className="app-frame">
      <Canvas/>
      </div>
    </div>
  );
}

export default App;
