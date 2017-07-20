import React, { Component } from 'react';
import {
  Form,
  Button,
  FormGroup,
  ControlLabel,
  FormControl
} from 'react-bootstrap';
import { connect } from 'react-redux';

import { actionInputSubmit, actionInputRoom, actionInputID } from './data/home';

/**
 * Validates a correct roomID.
 * @param {string} roomID 
 */
function validate(roomID) {
  if (roomID.length === 0) return false;
  if (isNaN(Number(roomID))) return false;
  const x = Number(roomID);
  if (x < 0) return false;
  return true;
}

function validateID(ID) {
  if (!(ID.length === 0 || ID.length === 32)) return false;
  return /^[a-zA-Z]*$/.test(ID);
}

export default connect(
  state => {
    return {
      enabled:
        validate(state.home.inputRoomID) && validateID(state.home.inputID),
      show: state.home.ws === null
    };
  },
  dispatch => {
    return {
      submit: () => dispatch(actionInputSubmit())
    };
  }
)(viewJumpTo);

// A jump to form to enter a custom room.
function viewJumpTo({ submit, enabled, show }) {
  if (!show) return null;
  return (
    <div className="text-center">
      <Form
        inline
        onSubmit={e => {
          e.preventDefault();
          submit();
        }}
      >
        <RoomInput />
        <IDInput style={{ marginLeft: '10px' }} />
        <Button
          type="submit"
          bsStyle="success"
          disabled={!enabled}
          style={{ marginLeft: '10px' }}
        >
          Go!
        </Button>
      </Form>
    </div>
  );
}

// An input for RoomID.
class viewRoomInput extends Component {
  handleChange(event) {
    this.props.change(event.target.value);
  }
  render() {
    const { roomID } = this.props;
    return (
      <FormGroup
        controlId="roomID"
        validationState={validate(roomID) ? 'success' : 'error'}
      >
        <ControlLabel>...or jump to room </ControlLabel>
        <FormControl
          type="text"
          placeholder="Enter the room ID"
          value={roomID}
          onChange={event => this.handleChange(event)}
          style={{ marginLeft: '10px' }}
        />
        <FormControl.Feedback />
      </FormGroup>
    );
  }
}

const RoomInput = connect(
  state => {
    return { roomID: state.home.inputRoomID };
  },
  dispatch => {
    return { change: val => dispatch(actionInputRoom(val)) };
  }
)(viewRoomInput);

// An input for ID.
class viewIDInput extends Component {
  handleChange(event) {
    this.props.change(event.target.value);
  }
  render() {
    const { ID } = this.props;
    return (
      <FormGroup
        controlId="ID"
        validationState={validateID(ID) ? 'success' : 'error'}
      >
        <ControlLabel>, maybe under the ID </ControlLabel>
        <FormControl
          type="text"
          placeholder="Your player ID (optional)"
          value={ID}
          onChange={event => this.handleChange(event)}
          style={{ marginLeft: '10px' }}
        />
        <FormControl.Feedback />
      </FormGroup>
    );
  }
}

const IDInput = connect(
  state => {
    return { ID: state.home.inputID };
  },
  dispatch => {
    return { change: val => dispatch(actionInputID(val)) };
  }
)(viewIDInput);
