import React, { Component } from 'react';
import countdown from 'countdown';
import { connect } from 'react-redux';
import { actionResponse } from './data/home';
import { Col, Button, Glyphicon, Row } from 'react-bootstrap';

import wavGameStart from './sounds/game_start.wav';
const gameStart = new Audio(wavGameStart);

// viewMatchmaker returns a matchmaking and ready check UI.
function viewMatchmaker({ found, changeDecision }) {
  if (found === null) {
    return null; // Nothing to see, no queuing.
  }
  if (found instanceof Date) {
    return <QueueClock since={found} />; // Expect the clock.
  }
  if (found === 'pending') {
    return <MatchFound changeDecision={changeDecision} />;
  }
  if (found === 'accepted') {
    return (
      <h4 className="text-center">
        <Glyphicon glyph="ok" /> You have accepted this match. Waiting for other
        players!
      </h4>
    );
  }
  if (found === 'rejected') {
    return (
      <h4 className="text-center">
        <Glyphicon glyph="remove" /> You have rejected this match. Waiting for
        other players...
      </h4>
    );
  }
}

export default connect(
  state => {
    return { found: state.home.found };
  },
  dispatch => {
    return {
      changeDecision: decision => dispatch(actionResponse(decision))
    };
  }
)(viewMatchmaker);

// QueueClock returns a clock that measures queue time.
class QueueClock extends Component {
  constructor() {
    super();
    this.state = { now: new Date() };
  }
  updateTime() {
    this.setState({ now: new Date() });
  }
  componentWillMount() {
    this.clock = setInterval(() => this.updateTime(), 1000);
  }
  componentWillUnmount() {
    clearInterval(this.clock);
  }
  render() {
    return (
      <h3 className="text-center">
        Queue Time: {countdown(this.state.now, this.props.since).toString()}
      </h3>
    );
  }
}

// MatchFound displays a box that asks for a ready check.
function MatchFound({ changeDecision }) {
  gameStart.play();
  return (
    <Row>
      <Col md={8} mdOffset={2} className="text-center">
        <h3>A match has been found!</h3>
        <Col md={5}>
          <Button
            bsStyle="success"
            bsSize="large"
            onClick={() => changeDecision(true)}
          >
            <Glyphicon glyph="ok" /> Accept
          </Button>
        </Col>
        <Col md={5} mdPush={2}>
          <Button
            bsStyle="danger"
            bsSize="large"
            onClick={() => changeDecision(false)}
          >
            <Glyphicon glyph="remove" /> Reject
          </Button>
        </Col>
      </Col>
    </Row>
  );
}
