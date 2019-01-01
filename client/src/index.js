import ReactDOM from "react-dom";
import React from "react";  // eslint-disable-line 
import Root from "./containers/Root";
import configureStore from "./store/configureStore";

import Popper from 'popper.js'
//require("bootstrap");
import "../node_modules/bootstrap/dist/css/bootstrap.css";


const store = configureStore();
ReactDOM.render(<Root store={store}/>, document.getElementById("root"));
