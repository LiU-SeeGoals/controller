import logo from './logo.svg';
import './App.css';
import React, { useState, useEffect, useRef } from "react";
import Canvas from "./components/Canvas/Canvas.js";
import Menu from "./components/Menu/Menu.js";

function App() {
  return (
    <div className="App">
      <Menu/>
      <Canvas/>
    </div>
  );
}

export default App;
