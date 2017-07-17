import React, { Component } from 'react';
import { connect } from 'react-redux';
import { actionConnect } from './data/game';

import countdown from 'countdown';
import YourTurn from './sounds/your_turn.wav';

const wavTurn = new Audio(YourTurn);

export default connect(state => {
  const g = state.game;
  return { myTurn: g.room.status[g.myID] === 'turn' };
})(viewStats);

function viewStats({ myTurn }) {
  return (
    <div>
      <h3>Game Status</h3>
      <hr />
      <div style={{ fontSize: '90%' }} className="text-center">
        <WSStatus />
        <TurnCount />
        {myTurn ? <TurnAnnouncer /> : null}
        <TurnClock />
      </div>
    </div>
  );
}

class viewWSStatus extends Component {
  componentWillMount() {
    if (this.props.ws === null) this.props.connect();
  }
  render() {
    const { ws } = this.props;
    if (ws === null)
      return (
        <h5>
          Websocket: <span className="text-danger">Not connected</span>
        </h5>
      );
    if (ws instanceof WebSocket)
      return (
        <h5>
          Websocket: <span className="text-success">Connected</span>
        </h5>
      );
    return null;
  }
}

const WSStatus = connect(
  state => {
    return { ws: state.game.ws };
  },
  dispatch => {
    return {
      connect: () => dispatch(actionConnect())
    };
  }
)(viewWSStatus);

// Shows the turn count.
function viewTurnCount({ count }) {
  return (
    <h5>
      Turns passed: <b>{count}</b>
    </h5>
  );
}

const TurnCount = connect(state => {
  return {
    count: state.game.room.sentences.filter(item => !item.system).length
  };
})(viewTurnCount);

// Displays how much time is left.
class viewTurnClock extends Component {
  constructor() {
    super();
    this.state = { now: new Date() };
  }
  componentWillMount() {
    this.timer = setInterval(() => {
      this.setState({ now: new Date() });
    }, 100);
  }
  componentWillUnmount() {
    clearInterval(this.timer);
  }
  render() {
    const { ended, turnStart, timeout, gameStart } = this.props;
    if (ended) return null;
    const tS = new Date(turnStart),
      gS = new Date(gameStart);
    if (tS.getTime() < gS.getTime()) {
      return <h5>Game starting soon...</h5>;
    }
    const turnEnd = new Date(new Date(turnStart).getTime() + timeout * 1000);
    if (turnEnd >= this.state.now) {
      return (
        <h4>
          {countdown(
            turnEnd,
            this.state.now,
            countdown.SECONDS,
            1,
            1
          ).toString()}
        </h4>
      );
    }
    return null;
  }
}

const TurnClock = connect(state => {
  const s = state.game;
  return {
    ended: s.ended !== null,
    gameStart: s.room.start,
    turnStart: s.room.current,
    timeout: s.room.timeout
  };
})(viewTurnClock);

// Announces it's your turn
function TurnAnnouncer() {
  wavTurn.play();
  return <h4>Your Turn!</h4>;
}
