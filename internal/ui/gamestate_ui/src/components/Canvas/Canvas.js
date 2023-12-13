import React, { useState, useEffect, useRef } from "react";
import "./Canvas.css";

function Canvas(props) {
    const canvasRef = useRef(null);
    const speed_arrow_color = 'rgba(0, 128, 128, 1)'
    let scaleFactor = 1;

    const canvasInit = (event) => {
        const canvas = canvasRef.current;
        // Adjust canvas width to match the image width
        canvas.width = event.target.width;
        canvas.height = event.target.height;

        scaleFactor = canvas.width / 6600; // Make the canvas independet on window scaling
        canvas.addEventListener('click', handleClick); // event for mouse click
        draw();
    };

    // draws everything on canvas
    function draw() {
        const canvas = canvasRef.current;
        const context = canvas.getContext('2d');

        // Clear the canvas
        context.clearRect(0, 0, canvas.width, canvas.height);
        
        // drawing all of gamestate here
        var gameState = props.gameState;
        gameState.map((robot) => {
            drawRobot(context, robot);
        });
    }

    const drawRobot = (context, robot) => {
        let x = robot.x;
        let y = robot.y;
        let radius = 180;
        let angle = Math.atan2(robot.speed_y, robot.speed_x);
        let arrowLength = 10* Math.sqrt(robot.speed_x * robot.speed_x + robot.speed_y * robot.speed_y);
        let arrowThickness = 3;
        let colorMap = {"yellow": "rgba(255, 255, 0, 1)", "blue": "rgba(0, 0, 255, 1)"};
        let color = colorMap[robot.team];

        drawArrow(context, x, y, angle, arrowLength, arrowThickness);
        drawCircle(context, x, y, radius, color);

        // Draw a black circle in the robot if it is selected
        if (robot.selected) {
            drawCircle(context, x, y, radius/3, 'rgba(0, 0, 0, 1)');
        }
    };

    const drawCircle = (context, x, y, radius, color) => {
        context.beginPath();
        context.arc(x*scaleFactor, y*scaleFactor, radius*scaleFactor, 0, 2 * Math.PI);
        context.strokeStyle = 'rgba(0, 0, 0, 0)'; // make the border transparent
        context.fillStyle = color;
        context.fill();
        context.stroke();
    };

    const drawArrow = (context, x, y, angle, length, thickness, color) => {
        context.beginPath();

        // Calculate the starting point of the arrow (on the circle)
        const startX = x;
        const startY = y;

        // Calculate the end point of the arrow
        const endX = x + length * Math.cos(angle);
        const endY = y + length * Math.sin(angle);

        // Draw the line for the arrow
        context.beginPath();
        context.moveTo(startX * scaleFactor, startY * scaleFactor);
        context.lineTo(endX * scaleFactor, endY * scaleFactor);
        context.strokeStyle = speed_arrow_color;
        context.lineWidth = thickness;
        context.stroke();

        // Draw the arrow head
        const headlen = 3;
        const angle1 = angle - Math.PI / 7;
        const angle2 = angle + Math.PI / 7;
        const headX = endX - headlen * Math.cos(angle1);
        const headY = endY - headlen * Math.sin(angle1);

        context.beginPath();
        context.moveTo(endX * scaleFactor, endY * scaleFactor);
        context.lineTo(headX * scaleFactor, headY * scaleFactor);
        context.lineTo((endX - headlen * Math.cos(angle2)) * scaleFactor, (endY - headlen * Math.sin(angle2)) * scaleFactor);
        context.lineTo(endX * scaleFactor, endY * scaleFactor);
        context.fillStyle = speed_arrow_color;
        context.fill();
        context.lineWidth = thickness;
        context.stroke();
    };

    const handleClick = (event) => {
        const rect = canvasRef.current.getBoundingClientRect();
        const x = (event.clientX - rect.left)/scaleFactor;
        const y = (event.clientY - rect.top)/scaleFactor;
        let gameStateCopy = [...props.gameState];

        // Check if clicked on a robot
        let clickedOnRobot = false;
        let clickedOnRobotIndex = -1;
        props.gameState.map((robot, index) => {
            if (Math.pow(x - robot.x, 2) + Math.pow(y - robot.y, 2) < Math.pow(3, 2)) {
                clickedOnRobot = true;
                clickedOnRobotIndex = index;
            }
        });  

        // If a robot is selected, move it here
        if (!clickedOnRobot) {
            props.gameState.map((robot, index) => {
                if (robot.selected) {
                    gameStateCopy[index]["x"] = x;
                    gameStateCopy[index]["y"] = y;
                }
            }); 
        }

        // deselect all robots
        gameStateCopy.map((robot, index) => {
            if(clickedOnRobot && index === clickedOnRobotIndex) {
                gameStateCopy[index]["selected"] = !gameStateCopy[index]["selected"];
            } else {
                gameStateCopy[index]["selected"] = false;
            }
        });
        props.setGameState(...gameStateCopy);
        draw();
    };

    return (
        <div className="canvasContainer">
            
            <div className="canvasBackground">
                <img src="./background.svg" alt="canvas" onLoad={canvasInit} />
            </div>
            <canvas className="myCanvas" ref={canvasRef} height={500} />
        </div>
    );
}

export default Canvas;
