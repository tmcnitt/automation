import React from 'react'
import {connect} from "react-redux";
import {getDevices} from '../actions';
import Device from './Device'

class Root extends React.Component {
	constructor(props){
		super(props)

		props.getDevices()
	}

	render(){
		if(this.props.devices != null && this.props.devices.length){
			return ( <p> Test </p>)
		}
		return (
			<div>
				{Object.keys(this.props.devices).map( (key, i) => {
					return <Device store={this.props.store} device={this.props.devices[key]} key={i} />;
				})}
			</div>
		)
	}
}

const mapStateToProps = (state) => ({ devices: state.devices.devices });

const mapDispatchToProps = (dispatch) => {
	return {
		getDevices: () => { dispatch(getDevices()); },
	};
};

export default connect(mapStateToProps, mapDispatchToProps)(Root);