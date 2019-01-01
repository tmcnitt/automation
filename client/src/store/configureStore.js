import {createStore, applyMiddleware, compose} from "redux";
import ReduxThunk from "redux-thunk";
import rootReducer from "../reducers";

const enhancer = compose(
	applyMiddleware(ReduxThunk),
	//window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
);

export default function configureStore(initialState) {
	// Note: only Redux >= 3.1.0 supports passing enhancer as third argument.
	// See https://github.com/rackt/redux/releases/tag/v3.1.0
	const store = createStore(rootReducer, initialState, enhancer);

	return store;
}
