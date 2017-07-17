// The "game_id" data struct.

/*
Contains
"gameID": <int | null>,
"ID": <string | null>,
*/

const ROOM_ID = 'ID_ROOM_ID_CHANGE';
export const ID = 'ID_CHANGE';
export const RESET = 'RESET';

/**
 * Changes game ID on action.
 * @param {Number | null} id 
 */
export function actionRoomID(id) {
  return {
    type: ROOM_ID,
    payload: id
  };
}

/**
 * Changes Room ID on action.
 * @param {Number?} roomID 
 * @param {{type: string, payload: any}} action
 * @return {Number?}
 */
export function reduceRoomID(roomID = null, action) {
  switch (action.type) {
    case RESET:
      return null;
    case ROOM_ID:
      return action.payload;
    default:
      return roomID;
  }
}

/**
 * Changes ID on action.
 * @param {string} id 
 */
export function actionID(id) {
  return {
    type: ID,
    payload: id
  };
}

/**
 * Resets the game.
 */
export function actionReset() {
  return { type: RESET };
}

/**
 * Changes ID on action.
 * @param {string} Id
 * @param {{type: string, payload: any}} action
 * @return {string}
 */
export function reduceID(Id = '', action) {
  switch (action.type) {
    case RESET:
      return '';
    case ID:
      return action.payload;
    default:
      return Id;
  }
}
