import React from 'react';
import { Row, Col, Image } from 'react-bootstrap';

import logo from './logo.png';

export default function Welcome() {
  return (
    <Row>
      <Col md={3}>
        <Image src={logo} responsive />
      </Col>
      <Col md={9}>
        <h2>Hakkero Project!</h2>
        <h4>
          A mini web novel game where players strives to win by writing the best
          collaborative story they can come up with.
        </h4>
        <div>
          The rules are simple:
          <ul>
            <li>You are given a topic, along with some other players</li>
            <li>
              The goal is to write the best story...<em> together</em>
            </li>
            <li>Each take turn and make one sentence in a specified time.</li>
            <li>You can skip your turn, but there's no turning back.</li>
            <li>The winner keeps the story alive until sunrise.</li>
          </ul>
          <b>
            Respect other players. If you have no idea, do not spam. Keep the
            game clear!
          </b>
        </div>
      </Col>
    </Row>
  );
}
