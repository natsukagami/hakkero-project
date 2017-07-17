import React from 'react';
import { Grid } from 'react-bootstrap';
import { connect } from 'react-redux';
import Home from './Home';
import Game from './Game';

function viewApp({ roomID, test }) {
  return (
    <Grid className="App" fliud={roomID !== null}>
      {roomID === null ? <Home /> : <Game />}
    </Grid>
  );
}

export default connect(
  state => {
    return { roomID: state.roomID };
  },
  dispatch => {
    return { test: () => dispatch({ type: 'ID_ROOM_ID_CHANGE', payload: 1 }) };
  }
)(viewApp);
