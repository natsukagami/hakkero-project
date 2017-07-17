import React, { Component } from 'react';
import { Row, Col, FormControl, Button } from 'react-bootstrap';
import { connect } from 'react-redux';

import { actionSubmit, actionSentence } from './data/game';

export default connect(state => {
  return { myTurn: state.game.room.status[state.game.myID] === 'turn' };
})(viewInput);

// The input section
function viewInput({ myTurn }) {
  return (
    <Row>
      <Col md={10}>
        <TextInput />
      </Col>
      <Col md={2}>
        {myTurn ? <Buttons /> : null}
      </Col>
    </Row>
  );
}

class viewTextInput extends Component {
  // Handle textarea change.
  handleChange(event) {
    const str = event.target.value;
    const len = str.length;
    if (str[len - 1] === '\n') {
      // Enter! Process this
      if (this.props.myTurn && str.length - 1 > 0) this.props.submit();
      return;
    }
    this.props.update(str);
  }
  render() {
    return (
      <FormControl
        componentClass="textarea"
        placeholder="Type your sentence here!"
        rows={4}
        onChange={e => this.handleChange(e)}
        value={this.props.sentence}
      />
    );
  }
}

const TextInput = connect(
  state => {
    return {
      sentence: state.game.sentence,
      myTurn: state.game.room.status[state.game.myID] === 'turn'
    };
  },
  dispatch => {
    return {
      submit: () => dispatch(actionSubmit(false)),
      update: str => dispatch(actionSentence(str))
    };
  }
)(viewTextInput);

function viewButtons({ submit, skip, canSubmit }) {
  return (
    <div>
      <Button
        className="form-control"
        bsStyle="success"
        disabled={!canSubmit}
        onClick={submit}
      >
        Go!
      </Button>
      <hr />
      <Button className="form-control" bsStyle="danger" onClick={skip}>
        Skip
      </Button>
    </div>
  );
}

const Buttons = connect(
  state => {
    return { canSubmit: state.game.sentence.length > 0 };
  },
  dispatch => {
    return {
      submit: () => dispatch(actionSubmit(false)),
      skip: () => dispatch(actionSubmit(true))
    };
  }
)(viewButtons);
