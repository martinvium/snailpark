import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import { connect } from 'react-redux'
import { receiveMessageAction, repeatingPingAction, startAction } from './actions'
import GameState from './components/GameState'
import CardList from './components/CardList'

const socketUrl = () => {
  const gameId = Date.now()
  const path = 'game/connect?gameId=' + gameId
  // const l = window.location
  // return ((l.protocol === "https:") ? "wss://" : "ws://") + l.host + l.pathname + path
  return 'ws://localhost:8081/' + path
}

class App extends Component {
  componentDidMount = () => {
    const { dispatch } = this.props;
    const socket = new WebSocket(socketUrl())

    socket.onmessage = (event) => receiveMessageAction(dispatch, event)
    socket.onopen = () => {
      dispatch(repeatingPingAction(socket))
      dispatch(startAction(socket))
    }
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p>

        <GameState/>
        <CardList/>
      </div>
    );
  }
}

const ConnectedApp = connect()(App)

export default ConnectedApp;
