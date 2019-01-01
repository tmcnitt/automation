import {createReducer} from "../utils";
import {GET_DEVICES_SUCCESS, GET_DEVICES_FAILURE, SET_DEVICE, SET_NAME_SUCESS, ADD_DEVICE_SUCCESS} from '../actions'
 
//TODO: How do init this
const initalState = {devices: {}};

const devices = createReducer(initalState, {
	[GET_DEVICES_SUCCESS](state,action){
		let devices = {}
		for (let i = 0; i < action.devices.length; i++){
			let device = action.devices[i]
			devices[device["sigs"]["zwave:nodeId"]] = device
		}
		
		return Object.assign({}, state, {
			devices: devices
		})
	},
	
});

export default devices;
