// The "game" data struct.

/* 
Contains
"game": {
	"ws": <WebSocket | error | null>,
	"room": {
		id: <int>,
		sentences: [
			{content: <string>, owner: <int>, system: <bool>}
		]
		members: <string[]>,
		status: <string[]>,
		start: Date,
		current: Date,
		timeout: int,
	} | null,
	"ended": <int | null>,
	"myID": <int, null>,
	"sentence": <string>,
}
*/

import Config from '../config';
import { RESET } from './id';
import { combineReducers } from 'redux';

// Events
const WS = 'GAME_WS_CHANGE';
const ROOM = 'GAME_ROOM_CHANGE';
const ROOM_SENTENCE = 'GAME_ROOM_SENTENCE';
const MYID = 'GAME_MYID_CHANGE';
const SENTENCE = 'GAME_SENTENCE_CHANGE';
const ENDED = 'GAME_ENDED_CHANGE';

/**
 * Handles incoming messages.
 * @param {WebSocket} ws 
 * @param {MessageEvent} event 
 * @param {Function} dispatch 
 */
function gameMessageDispatcher(ws, event, dispatch, state) {
  try {
    console.log(event.data);
    const data = JSON.parse(event.data);
    const payload = data.message;
    switch (data.type) {
      case 'index':
        return dispatch({ type: MYID, payload: payload.index });
      case 'turn':
        return dispatch({ type: ROOM, payload: payload });
      case 'sentence':
        return dispatch({ type: ROOM_SENTENCE, payload: payload });
      case 'end':
        return dispatch({ type: ENDED, payload: payload.Winner });
      default:
    }
  } catch (e) {
    console.log(e);
  }
}

/**
 * Fetches room data from POST.
 */
export function initRoom() {
  return async (dispatch, getState) => {
    let res = null;
    try {
      const id = getState().roomID;
      res = await fetch(Config.server + `/rooms/${id}`, { method: 'POST' });
      res = await res.json();
      res.timeout /= 1000000000;
    } catch (e) {
      console.log(e);
      dispatch({
        type: ROOM,
        payload: new Error(e)
      });
    }
    if (res !== null)
      dispatch({
        type: ROOM,
        payload: res
      });
  };
}

/**
 * Connects to the WebSocket server.
 */
export function actionConnect() {
  return (dispatch, getState) => {
    const state = getState();
    const { roomID, ID } = state;
    const ws = new WebSocket(Config.wsServer + `/rooms/${roomID}?player=${ID}`);
    dispatch({ type: WS, payload: ws });
    ws.addEventListener('error', ev => {
      // Error occurred
      dispatch({ type: WS, payload: new Error(ev.toString()) });
    });
    ws.addEventListener('close', ev => {
      // Disconnected
      dispatch({ type: WS, payload: null });
    });
    ws.addEventListener('message', ev =>
      gameMessageDispatcher(ws, ev, dispatch, getState())
    );
  };
}

/**
 *  Changes sentence on action.
 * @param {string} sentence 
 */
export function actionSentence(sentence) {
  return {
    type: SENTENCE,
    payload: sentence
  };
}

/**
 * Submits the sentence.
 * @param {Boolean} skip
 */
export function actionSubmit(skip) {
  return (dispatch, getState) => {
    const state = getState();
    const ws = state.game.ws;
    const sentence = state.game.sentence;
    if (ws === null || ws instanceof Error)
      throw new Error('Unexpected submit.'); // No use.
    ws.send(JSON.stringify({ skip: skip, content: sentence }));
    dispatch({
      type: SENTENCE,
      payload: ''
    });
  };
}

/**
 * Reduces WS object.
 * @param {WebSocket | Error} ws 
 * @param {Action} action
 * @return {WebSocket | Error}
 */
function reduceWS(ws = null, action) {
  if (action.type === WS) {
    return action.payload;
  }
  return ws;
}

/**
 * Reduce the array of sentences.
 * @param {Array<any>} sentences 
 * @param {*} action
 * @return {Array<any>}
 */
function reduceRoomSentences(sentences = [], action) {
  const nw = sentences.slice();
  if (action.type === ROOM_SENTENCE) {
    nw.splice(action.payload.pos, 0, action.payload.sentence);
    return nw;
  }
  return sentences;
}

/**
 * Reduces Room object.
 * @param {any} room 
 * @param {Action} action
 * @return {any}
 */
function reduceRoom(room = null, action) {
  switch (action.type) {
    case ROOM:
      if (action.payload === null || action.payload instanceof Error)
        return action.payload;
      return Object.assign({}, room || {}, action.payload);
    case ROOM_SENTENCE:
      return Object.assign({}, room, {
        sentences: reduceRoomSentences(room.sentences, action)
      });
    default:
      return room;
  }
}

/**
 * Reduces "ended" status.
 * @param {Number} ended 
 * @param {*} action
 * @return {Number}
 */
function reduceEnded(ended = null, action) {
  if (action.type === ENDED) {
    return action.payload;
  }
  return ended;
}

/**
 * Reduces myID indicator.
 * @param {Number} myID 
 * @param {*} action
 * @return {Number}
 */
function reduceMyID(myID = null, action) {
  if (action.type === MYID) {
    return action.payload;
  }
  return myID;
}

/**
 * Reduces current sentence.
 * @param {string} sent 
 * @param {*} action
 * @return {string}
 */
function reduceSentence(sent = '', action) {
  if (action.type === SENTENCE) {
    return action.payload;
  }
  return sent;
}

export function reduceGame(game = {}, action) {
  const rG = combineReducers({
    ws: reduceWS,
    room: reduceRoom,
    ended: reduceEnded,
    myID: reduceMyID,
    sentence: reduceSentence
  });
  if (action.type === RESET) {
    return rG({}, action);
  }
  return rG(game, action);
}
