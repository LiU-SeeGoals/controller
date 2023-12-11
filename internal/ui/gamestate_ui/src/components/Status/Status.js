import React, { useState, useEffect, useRef } from "react";
import "./Status.css";

function Status(props) {
    const [selectedMenu, setSelectedMenu] = useState(1);

    function selectMenu(menuId) {
        setSelectedMenu(menuId);
      }
    
    let robotStatus = [
        {id: 0, x: 50, y: 50, angle: 3, action: "idle", battery: 30, speed: 10},
        {id: 1, x: 50, y: 50, angle: 3, action: "idle", battery: 30, speed: 10},
        {id: 2, x: 50, y: 50, angle: 3, action: "idle", battery: 30, speed: 10},
        {id: 3, x: 50, y: 50, angle: 3, action: "idle", battery: 30, speed: 10},
        {id: 4, x: 50, y: 50, angle: 3, action: "idle", battery: 30, speed: 10},
        {id: 5, x: 50, y: 50, angle: 3, action: "idle", battery: 30, speed: 10}]

    return (
        <div className="statusContainer">
            {robotStatus.map((robot) => (
            <div key={robot.id} className="robotCard">
            <p><b>Robot ID: {robot.id}</b></p>
            <p>X: {robot.x}</p>
            <p>Y: {robot.y}</p>
            <p>Angle: {robot.angle}</p>
            <p>Action: {robot.action}</p>
            <p>Battery: {robot.battery}%</p>
            <p>Speed: {robot.speed}</p>
        </div>
      ))}
        </div>
    );
}

export default Status;
