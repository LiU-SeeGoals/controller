import React, { useState, useEffect, useRef } from "react";
import "./Canvas.css";

function Canvas(props) {
    const canvasRef = useRef(null);
    const speed_arrow_color = 'rgba(0, 128, 128, 1)'
    let scaleFactor = 1;

    function draw() {
        const canvas = canvasRef.current;
        const context = canvas.getContext('2d');

        // drawing all of gamestate here
        var gameState = props.gameState;
        gameState.map((robot) => {
            drawRobot(context, robot);
        
        });
        // // Your drawing code goes here
        // context.fillStyle = 'rgba(255, 0, 0, 0.5)';
        // context.fillRect(0, 0, canvas.width, canvas.height);
        
        // drawRobot(context, 50, 50, 3, 3.14, 30);
        // drawCircle(context, 75, 50, 3);
      }
    
    const drawCircle = (context, x, y, radius) => {
        context.beginPath();
        context.arc(x*scaleFactor, y*scaleFactor, radius*scaleFactor, 0, 2 * Math.PI);
        context.fillStyle = 'rgba(255, 0, 0, 1)';
        context.fill();
        context.stroke();
    };

    const drawRobot = (context, robot) => {
        let x = robot.x;
        let y = robot.y;
        let radius = 3;
        let angle = Math.atan2(robot.speed_y, robot.speed_x);
        let arrowLength = 10* Math.sqrt(robot.speed_x * robot.speed_x + robot.speed_y * robot.speed_y);
        let arrowThickness = 3;
        let colorMap = {"yellow": "rgba(255, 255, 0, 1)", "blue": "rgba(0, 0, 255, 1)"};
        let color = colorMap[robot.team];
        drawArrow(context, x, y, angle, arrowLength, arrowThickness);
        // Draw the circle
        context.beginPath();
        context.arc(x * scaleFactor, y * scaleFactor, radius * scaleFactor, 0, 2 * Math.PI);
        context.strokeStyle = 'rgba(255, 0, 0, 0)';
        context.fillStyle = color;
        context.fill();
        context.stroke();

        if (robot.selected) {
            context.beginPath();
            context.arc(x * scaleFactor, y * scaleFactor, radius/5 * scaleFactor, 0, 2 * Math.PI);
            context.strokeStyle = 'rgba(0, 0, 0, 0)';
            context.fillStyle = 'rgba(0, 0, 0, 1)';
            context.fill();
            context.stroke();
        }

    
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
        
        // deselect all robots
        gameStateCopy.map((robot, index) => {
            gameStateCopy[index]["selected"] = false;
        });

        props.gameState.map((robot, index) => {
            if (Math.pow(x - robot.x, 2) + Math.pow(y - robot.y, 2) < Math.pow(3, 2)) {
                console.log("selected: " + index);
                gameStateCopy[index]["selected"] = !gameStateCopy[index]["selected"];
                props.setGameState(...gameStateCopy);
                console.log(props.gameState)
                draw();
            }
        });  
    };

    // const handleMouseMove = (event) => {
    //     const context = canvas.getContext('2d');
    //     context.clearRect(0, 0, canvas.width, canvas.height); // Clear the canvas to remove previous circles.
    //     draw();
        
    //     const x = event.clientX - canvas.getBoundingClientRect().left;
    //     const y = event.clientY - canvas.getBoundingClientRect().top;
        
    //     context.beginPath();
    //     context.arc(x * scaleFactor, y * scaleFactor, radius * scaleFactor, 0, 2 * Math.PI);
    //     context.strokeStyle = 'rgba(0, 0, 0, 0)';
    //     context.fillStyle = 'rgba(0, 0, 0, 1)';
    //     context.fill();
    //     context.stroke();
    //   };

    const handleImageLoad = (event) => {
        const canvas = canvasRef.current;
        console.log("canvas width: " + event.target.width); 
        // Adjust canvas width to match the image width
        canvas.width = event.target.width;
        canvas.height = event.target.height;
        // Make the canvas independet on window scaling
        scaleFactor = canvas.width / 100; 
        canvas.addEventListener('click', handleClick);
        // canvas.addEventListener('mousemove', handleMouseMove);
        draw();
      };

    return (
        <div className="canvasContainer">
            
            <div className="canvasBackground">
                <img src="./background.png" alt="canvas" onLoad={handleImageLoad} />
            </div>
            <canvas className="myCanvas" ref={canvasRef} height={500} />
        </div>
    );
}

export default Canvas;
