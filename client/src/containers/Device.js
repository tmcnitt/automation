import React from 'react'
import {connect} from "react-redux";
import {setDevice} from '../actions';

class Device extends React.Component {
	render(){
		return (
			<div>
				<p>{this.props.device.name || "Name"}</p>
				{Object.keys(this.props.device.status).map( (key, i) => {
					return <p>{key} : <input type="text" ref={key}  placeholder={(this.props.device.status[key]).toString()} /> <input type="button" onClick={() => { this.set(key) }} /></p>
				})}
			</div>
		)
	}

	set(key){
		this.props.setDevice(this.props.device.ID, key, this.refs[key].value)
	}
}

const mapStateToProps = (state) => ({});

const mapDispatchToProps = (dispatch) => {
	return {
		setDevice: (id, cc, status) => { dispatch(setDevice(id,cc,status)); },
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Device);