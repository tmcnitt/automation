import {combineReducers} from "redux";
import {routerReducer} from "react-router-redux";

import devices from "./devices";

export default combineReducers({
	routing: routerReducer,
	devices,
});
