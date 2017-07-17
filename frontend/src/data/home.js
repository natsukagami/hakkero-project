// The "home" data struct.

/*
Contains
"home": {
	"username": <string>,
	"ID": <string |  null>,
	"ws": <WebSocket | null | Error>,
	"console": [<string>, ...],
	"queueSize": <int>,
	"found": <Date | "pending" | "accepted" | "rejected" | null>,
	"roomID": <int | null>
}
*/

import { combineReducers } from 'redux';
import { actionRoomID, actionID, ID, RESET } from './id';
import Config from '../config';

const USERNAME = 'HOME_USERNAME_CHANGE';
const WS = 'HOME_WS_CHANGE';
const CONSOLE = 'HOME_CONSOLE_CHANGE';
const QUEUE_SIZE = 'HOME_QUEUE_SIZE_CHANGE';
const FOUND = 'HOME_FOUND_CHANGE';

/**
 * Changes the current username.
 * @param {string} username 
 */
export function actionUsername(username) {
  console.log(username);
  return {
    type: USERNAME,
    payload: username
  };
}

/**
 * Dispatches actions according to event.
 * @param {WebSocket} w
 * @param {MessageEvent} event
 * @param {Function} dispatch 
 */
function queueMessageDispatcher(w, event, dispatch) {
  try {
    const data = JSON.parse(event.data);
    const payload = data.message;
    switch (data.type) {
      case 'announcement':
        if (payload.success) {
          setTimeout(() => dispatch(actionRoomID(payload.room)), 1000);
        } else dispatch({ type: FOUND, payload: new Date() });
        return dispatch({ type: CONSOLE, payload: payload.announcement });
      case 'size':
        return dispatch({ type: QUEUE_SIZE, payload: payload.size });
      case 'found':
        return dispatch({ type: FOUND, payload: 'pending' });
      case 'ID':
        return dispatch(actionID(payload.ID));
      default:
    }
  } catch (e) {
    console.log(e);
  }
}

/**
 * Connects to the websocket.
 * @param {string} username
 */
export function actionConnect() {
  return (dispatch, getState) => {
    const ws = new WebSocket(
      Config.wsServer + '/queue?username=' + getState().home.username
    );
    ws.addEventListener('open', ev => {
      dispatch({
        type: WS,
        payload: ws
      });
      dispatch({
        type: FOUND,
        payload: new Date()
      });
    });
    ws.addEventListener('error', ev => {
      dispatch({
        type: WS,
        payload: new Error(require('util').inspect(ev))
      });
    });
    ws.addEventListener('close', ev => {
      dispatch({
        type: WS,
        payload: null
      });
      dispatch({
        type: FOUND,
        payload: null
      });
    });
    ws.addEventListener('message', msg =>
      queueMessageDispatcher(ws, msg, dispatch)
    );
  };
}

/**
 * Accept or reject the ready check.
 * @param {Boolean} accept 
 */
export function actionResponse(accept) {
  return (dispatch, getState) => {
    const state = getState();
    const ws = state.home.ws;
    if (ws === null || ws instanceof Error)
      throw new Error('Unexpected response.'); // No use.
    ws.send(JSON.stringify({ accepted: accept }));
    dispatch({
      type: FOUND,
      payload: accept ? 'accepted' : 'rejected'
    });
  };
}

/**
 * Adds a line to the console.
 * @param {string} text 
 */
export function actionConsole(text) {
  return {
    type: CONSOLE,
    payload: text
  };
}

/**
 * Returns a welcome message.
 */
export function actionWelcome() {
  return async (dispatch, getState) => {
    let res = await fetch(Config.server, { method: 'POST' });
    res = await res.json();
    dispatch(actionConsole(res.toString()));
  };
}

/**
 * Changes username on action.
 * @param {string} username 
 * @param {{type: string, payload: any}} action
 * @return {string}
 */
function reduceUsername(username = '', action) {
  switch (action.type) {
    case USERNAME:
      return action.payload;
    default:
      return username;
  }
}

/**
 * Changes the WebSocket instance on action.
 * @param {(WebSocket | Error)?} ws 
 * @param {{type: string, payload: any}} action
 * @return {(WebSocket | Error)?}
 */
function reduceWs(ws = null, action) {
  switch (action.type) {
    case WS:
      return action.payload;
    default:
      return ws;
  }
}

/**
 * Logs every action on the console.
 * @param {string[]} console 
 * @param {{type: string, payload: any}} action
 * @return {string[]}
 */
function reduceConsole(cs = [], action) {
  const ncs = cs.slice();
  let sentence = null;
  switch (action.type) {
    case ID:
      sentence = `Your ID is ${action.payload}`;
      break;
    case WS:
      if (action.payload === null)
        sentence = 'You have been disconneced from queue.';
      else if (action.payload instanceof WebSocket)
        sentence = 'You are now connected to the server! Currently in queue...';
      else sentence = 'An error has occured: ' + action.payload.toString();
      break;
    case CONSOLE:
      sentence = action.payload;
      break;
    case FOUND:
      switch (action.payload) {
        case 'pending':
          sentence =
            'A room has been found! Please accept or reject the ready check.';
          break;
        case 'accepted':
          sentence = 'You have accepted the match!';
          break;
        case 'rejected':
          sentence = 'You have rejected the match.';
          break;
        default:
      }
      break;
    default:
  }
  if (sentence !== null) {
    ncs.unshift(new Date().toString() + ': ' + sentence);
    return ncs;
  }
  return cs;
}

/**
 * Changes queue size on action.
 * @param {Number} queueSize 
 * @param {{type: string, payload: any}} action
 * @return {Number}
 */
function reduceQueueSize(queueSize = 0, action) {
  switch (action.type) {
    case QUEUE_SIZE:
      return action.payload;
    default:
      return queueSize;
  }
}

/**
 * Changes found status on action.
 * @param {Date | "pending" | "accepted" | "rejected" | null} found 
 * @param {{type: string, payload: any}} action
 * @return {Date | "pending" | "accepted" | "rejected" | null}
 */
function reduceFound(found = null, action) {
  switch (action.type) {
    case FOUND:
      return action.payload;
    default:
      return found;
  }
}

export function reduceHome(home = {}, action) {
  const rH = combineReducers({
    username: reduceUsername,
    ws: reduceWs,
    console: reduceConsole,
    queueSize: reduceQueueSize,
    found: reduceFound
  });
  if (action.type === RESET) {
    return rH({ username: home.username }, action);
  }
  return rH(home, action);
}
