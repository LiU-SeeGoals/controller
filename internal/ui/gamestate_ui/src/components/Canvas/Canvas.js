import React, { useState, useEffect, useRef } from "react";
import "./Canvas.css";

function Canvas(props) {
    const canvasRef = useRef(null);
    const speed_arrow_color = 'rgba(0, 128, 128, 1)'
    let scaleFactor = 1;

    function draw() {
        const canvas = canvasRef.current;
        const context = canvas.getContext('2d');
    
        // // Your drawing code goes here
        // context.fillStyle = 'rgba(255, 0, 0, 0.5)';
        // context.fillRect(0, 0, canvas.width, canvas.height);
        
        drawRobot(context, 50, 50, 3, 30, 30);
        drawCircle(context, 75, 50, 3);
      }
    
    const drawCircle = (context, x, y, radius) => {
        context.beginPath();
        context.arc(x*scaleFactor, y*scaleFactor, radius*scaleFactor, 0, 2 * Math.PI);
        context.fillStyle = 'rgba(255, 0, 0, 1)';
        context.fill();
        context.stroke();
    };

    const drawRobot = (context, x, y, radius, angle, arrowLength, arrowThickness) => {

        drawArrow(context, x, y, angle, arrowLength, arrowThickness);
        // Draw the circle
        context.beginPath();
        context.arc(x * scaleFactor, y * scaleFactor, radius * scaleFactor, 0, 2 * Math.PI);
        context.strokeStyle = 'rgba(255, 0, 0, 0)';
        context.fillStyle = 'rgba(255, 0, 0, 1)';
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

    const handleImageLoad = (event) => {
        const canvas = canvasRef.current;
        console.log("canvas width: " + event.target.width); 
        // Adjust canvas width to match the image width
        canvas.width = event.target.width;
        canvas.height = event.target.height;
        // Make the canvas independet on window scaling
        scaleFactor = canvas.width / 100; 
    
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
