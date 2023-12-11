import React, { useState, useEffect, useRef } from "react";
import "./Menu.css";
import Status from "../Status/Status.js";
import Overview from "../Overview/Overview.js";
import Control from "../Control/Control.js";


function Menu(props) {
    const [selectedMenu, setSelectedMenu] = useState(1);

    function selectMenu(menuId) {
        setSelectedMenu(menuId);
      }
    

    return (
        <div className="menuContainer">
            <div className="headerMenu" >
                <div id={selectedMenu == 1 ? "selected" : "_"} onClick={() => selectMenu(1)}>Overview</div>
                <div id={selectedMenu == 2 ? "selected" : "_"} onClick={() => selectMenu(2)}>Status</div>
                <div id={selectedMenu == 3 ? "selected" : "_"} onClick={() => selectMenu(3)}>Control</div>
            </div>
            <div className="contentMenu">
                <div style={{ display: selectedMenu === 1 ? 'block' : 'none' }}><Overview/></div>
                <div style={{ display: selectedMenu === 2 ? 'block' : 'none' }}><Status/></div>
                <div style={{ display: selectedMenu === 3 ? 'block' : 'none' }}><Control/></div>
      
            </div>
        </div>
    );
}

export default Menu;
