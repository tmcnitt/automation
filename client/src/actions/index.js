export const API_ROOT = "http://" + window.location.hostname + ":80";

export const GET_DEVICES_SUCCESS = "GET_DEVICES_SUCCESS"
export const GET_DEVICES_FAILURE = "GET_DEVICES_FAILURE"

export const GET_SCENES_SUCCESS = "GET_SCENES_SUCCESS"
export const GET_SCENES_FAILURE = "GET_SCENES_FAILURE"

export const CREATE_SCENE_SUCCESS = "CREATE_SCENE_SUCCESS"
export const CREATE_SCENE_FAILURE = "CERATE_SCENE_FAILURE"

export const DELETE_SCENE_SUCCESS = "DELETE_SCENE_SUCESS"
export const DELETE_SCENE_FAILURE = "DELETE_SCENE_FAILURE"

export const SCHEDULE_SCENE_SUCCESS = "SCHEDULE_SCENE_SUCUESS"	
export const SCHEDULE_SCENE_FAILURE = "SCHEDULE_SCENE_FAILURE"

export const NAME_DEVICE_SUCCESS = "NAME_DEVICE_SUCCESS"
export const NAME_DEVICE_FAILURE = "NAME_DEVICE_FAILURE"

export const ADD_DEVICE_SUCCESS = "ADD_DEVICE_SUCCESS"
export const ADD_DEVICE_FAILURE = "ADD_DEVICE_FAILURE"

export function getDevices(){
	return function(dispatch){
		fetch(API_ROOT + "/devices", {
			method: "GET",
		}).then(response => {
			response.json().then(r => {
				if (response.ok){

					dispatch({
						type: GET_DEVICES_SUCCESS,
						devices: r
					})
				} else {
					dispatch({type: GET_DEVICES_FAILURE})
				}
			})
		})
	}
}

export function getScenes(){
	return function(dispatch){
		fetch(API_ROOT + "/scenes", {
			method: "GET",
		}).then(response => {
			if (response.ok){

				dispatch({
					type: GET_SCENES_SUCCESS,
					scenes: r
				});
			} else {
				dispatch({type: GET_SCENES_FAILURE});
			}
		})
	}
}

export function createScene(scene){
	return function(dispatch){
		fetch(API_ROOT + "/scenes", {
			method: "POST",
			body: JSON.stringify(scene)
		}).then(response => {
			response.json().then(r => {
				if(response.ok){
					scene.id = r.id

					dispatch({
						type: CREATE_SCENE_SUCCESS,
						scene: scene
					});
				} else {
					dispatch({type: CREATE_SCENE_FAILURE});
				}
			})
		})
	}
}

export function deleteScene(id){
	return function(dispatch){
		fetch(API_ROOT + `/scenes/${id}`, {
			method: "DELETE",
		}).then(response => {
			response.json().then(r => {
				if(response.ok){
					dispatch({type: DELETE_SCENE_SUCCESS, id: id});
				} else {
					dispatch({type: DELETE_SCENE_FAILURE});
				}
			})
		})
	}
}

export function scheduleScene(id, time){
	return function(dispatch){
		fetch(API_ROOT + `/scenes/${id}/schedule`, {
			method: "POST",
			body: JSON.stringify({time: time})
		}).then(response => {
			response.json().then(r => {
				if(response.ok){
					dispatch({type: SCHEDULE_SCENE_SUCCESS, time: time})
				} else {
					dispatch({type: SCHEDULE_SCENE_FAILURE})
				}
			})
		})
	}
}

//Update is handled by the websocket and not here
export function setDevice(id, cc, status){
	return function(dispatch){
		fetch(API_ROOT +`/set/${id}/${cc}/${status}`)
	}
}

export function addDevice(){
	return function(dispatch){
		fetch(API_ROOT + "/add").then(response => {
			response.json(r => {
				if (response.ok){
					dispatch({type: ADD_DEVICE_SUCCESS})
				}else {
					dispatch({type: ADD_DEVICE_FAILURE})
				}
			})
		})
	}
}

export function nameDevice(id, name){
	return function(dispatch){
		fetch(API_ROOT + "/name", {
			method: "POST",
			body: JSON.stringify({ID: id, name: name})
		}).then(response => {
			response.json(r => {
				if(response.ok){
					dispatch({type: NAME_DEVICE_SUCCESS})
				} else {
					dispatch({type: NAME_DEVICE_FAILURE})
				}
			})
		})
	}
}


