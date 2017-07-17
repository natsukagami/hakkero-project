import React, { Component } from 'react';
import AutoScroll from 'react-auto-scroll';
import { Well } from 'react-bootstrap';
import { connect } from 'react-redux';

export default connect(state => {
  return {
    sentences: state.game.room.sentences,
    players: state.game.room.members
  };
})(viewSentences);

// Display all sentences in a nice view.
function viewSentences({ sentences, players }) {
  return (
    <div>
      <h3 className="text-center">
        {sentences[0].content}
      </h3>
      <hr />
      <SentencesList sentences={sentences.slice(1)} players={players} />
    </div>
  );
}

class sentencesList extends Component {
  render() {
    const { sentences, players } = this.props;
    return (
      <Well
        style={{
          height: '40vh',
          objectFit: 'cover',
          overflowY: 'scroll',
          fontFamily: '"Courier New", Courier, monospace'
        }}
      >
        {sentences.map((item, id) => {
          if (item.system)
            return (
              <div className="text-right" key={id}>
                <em>
                  {item.content}
                </em>
              </div>
            );
          return (
            <div key={id} title={`written by ${players[item.owner]}`}>
              {item.content}
            </div>
          );
        })}
        {sentences.length === 0 ? 'Your sentences will be here...' : null}
      </Well>
    );
  }
}

const SentencesList = AutoScroll({
  property: 'sentences'
})(sentencesList);
