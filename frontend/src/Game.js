import React, { Component } from 'react';
import { Image, Row, Col } from 'react-bootstrap';
import { connect } from 'react-redux';
import { actionReset } from './data/id';
import { initRoom } from './data/game';

import Welcome from './Welcome';
import Loading from './loading.svg';
import Logo from './logo_without_name.png';

import PlayerList from './PlayerList';
import Sentences from './Sentences';
import Stats from './Stats';
import Input from './Input';

class viewGame extends Component {
  componentWillMount() {
    if (this.props.room === null) this.props.init();
  }
  render() {
    const { room, ws, game } = this.props;
    if (room === null) {
      return (
        <div>
          <h2 className="text-center">
            <Image width="100px" height="100px" circle src={Loading} />
            Loading...
          </h2>
          <Welcome />
          <div>
            <ResetLink>Return to home page</ResetLink>
          </div>
        </div>
      );
    }
    if (room instanceof Error) {
      return (
        <h2 className="text-center">
          An error occured while loading room data. Please{' '}
          <ResetLink>return to home page</ResetLink>.
        </h2>
      );
    }
    if (ws instanceof Error) {
      return (
        <h2 className="text-center">
          An error occured while connecting to server. Please{' '}
          <ResetLink>return to home page</ResetLink>.
        </h2>
      );
    }
    return (
      <div>
        <Title />
        <Row>
          <Col md={3}>
            <PlayerList />
          </Col>
          <Col md={7}>
            <Sentences />
          </Col>
          <Col md={2}>
            <Stats />
          </Col>
        </Row>
        <Input />
        {require('util').inspect(game)}
      </div>
    );
  }
}

export default connect(
  state => {
    return { room: state.game.room, ws: state.game.ws, game: state.game };
  },
  dispatch => {
    return {
      init: () => dispatch(initRoom())
    };
  }
)(viewGame);

// Creates a link that resets the page.
function viewResetLink({ children, reset }) {
  return (
    <a href="#reset" onClick={reset}>
      {children}
    </a>
  );
}

const ResetLink = connect(undefined, dispatch => {
  return { reset: () => dispatch(actionReset()) };
})(viewResetLink);

// Show a title.
function viewTitle({ loading, ended, roomID }) {
  let status = 'running';
  if (loading) status = 'loading';
  if (ended) status = 'ended';
  return (
    <Row>
      <Col md={1} xs={3}>
        <ResetLink>
          <Image src={Logo} responsive />
        </ResetLink>
      </Col>
      <Col md={11} xs={9}>
        <h1>
          Room {roomID} - Game {status}!
        </h1>
      </Col>
    </Row>
  );
}

const Title = connect(state => {
  return {
    loading: !(state.game.ws instanceof WebSocket),
    ended: state.game.ended !== null,
    roomID: state.roomID
  };
})(viewTitle);
