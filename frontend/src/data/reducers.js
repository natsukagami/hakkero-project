import { combineReducers } from 'redux';
import { reduceHome } from './home';
import { reduceRoomID, reduceID } from './id';
import { reduceGame } from './game';

export default combineReducers({
  home: reduceHome,
  roomID: reduceRoomID,
  ID: reduceID,
  game: reduceGame
});
