import React from 'react';
import { ListGroup, ListGroupItem, Badge } from 'react-bootstrap';
import { connect } from 'react-redux';

function viewPlayerList({ players, statuses, myID, ended }) {
  return (
    <div>
      <h3>Players</h3>
      <hr />
      <ListGroup>
        {players.map((name, id) => {
          const status = ended === id ? 'winner' : statuses[id];
          const style = statusStyle(status);
          return (
            <ListGroupItem key={id} bsStyle={style}>
              {id === myID
                ? <b>
                    {name}
                  </b>
                : name}
              <Badge pullRight>
                {status}
              </Badge>
            </ListGroupItem>
          );
        })}
      </ListGroup>
    </div>
  );
}

export default connect(state => {
  return {
    players: state.game.room.members,
    statuses: state.game.room.status,
    myID: state.game.myID,
    ended: state.game.ended
  };
})(viewPlayerList);

// Returns the status label of the player.
function statusStyle(status) {
  switch (status) {
    case 'winner':
    case 'active':
      return 'success';
    case 'turn':
      return 'info';
    case 'skipped':
      return 'danger';
    case 'disconnected':
      return 'warning';
    default:
      return '';
  }
}
