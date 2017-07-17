import React, { Component } from 'react';
import { connect } from 'react-redux';
import {
  PageHeader,
  Button,
  Form,
  FormControl,
  FormGroup,
  ControlLabel
} from 'react-bootstrap';
import { actionConnect, actionUsername, actionWelcome } from './data/home';

import JumpTo from './JumpTo';
import Matchmaker from './Matchmake';
import Welcome from './Welcome';

// Home is the front page of the app.
class Home extends Component {
  // Changes the title of the page.
  componentWillMount() {
    document.title = 'Home - Hakkero Project!';
  }
  render() {
    return (
      <div>
        <Title />
        <QueueForm />
        <Matchmaker />
        <br />
        <Console />
        <br />
        <JumpTo />
        <hr />
        <Welcome />
      </div>
    );
  }
}

export default Home;

// Title returns the animated Title on the front page.
function Title() {
  return (
    <PageHeader>
      <div className="text-center">
        Hakkero<sup>[dev]</sup>!
      </div>
    </PageHeader>
  );
}

// The view for the console.
class viewConsole extends Component {
  componentWillMount() {
    this.props.welcome();
  }
  render() {
    return (
      <pre>
        {this.props.text.map(line => line + '\n')}
      </pre>
    );
  }
}

// The Console module, that displays the logs.
const Console = connect(
  state => {
    return { text: state.home.console };
  },
  dispatch => {
    return { welcome: text => dispatch(actionWelcome()) };
  }
)(viewConsole);

// queueButton
function QueueButton({ enabled }) {
  return (
    <Button
      type="submit"
      bsStyle="success"
      disabled={!enabled}
      style={{ marginLeft: '10px' }}
    >
      Start!
    </Button>
  );
}

// Runs validation on username.
function validate(username) {
  return username.length === 0 || username.length > 20 ? 'error' : 'success';
}

class UsernameInput extends Component {
  // Handles username change
  handleChange(event) {
    this.props.changeUsername(event.target.value);
  }
  render() {
    return (
      <FormGroup
        controlId="username"
        validationState={validate(this.props.username)}
      >
        <ControlLabel>Your name: </ControlLabel>
        <FormControl
          type="text"
          placeholder="Enter your name"
          value={this.props.username}
          onChange={event => this.handleChange(event)}
          readOnly={!this.props.enabled}
          style={{ marginLeft: '10px' }}
        />
        <FormControl.Feedback />
      </FormGroup>
    );
  }
}

// The form renderer.
function viewForm({ username, enabled, submit, changeUsername }) {
  return (
    <div className="text-center">
      <Form onSubmit={submit} inline>
        <UsernameInput
          username={username}
          changeUsername={changeUsername}
          enabled={enabled}
        />
        <QueueButton enabled={enabled && validate(username) === 'success'} />
      </Form>
    </div>
  );
}

const QueueForm = connect(
  state => {
    return {
      username: state.home.username,
      enabled: state.home.ws === null
    };
  },
  dispatch => {
    return {
      submit: e => {
        e.preventDefault();
        return dispatch(actionConnect());
      },
      changeUsername: username => dispatch(actionUsername(username))
    };
  }
)(viewForm);
