import React from 'react';
import { connect } from 'react-redux';

let GameState = ({ gameState }) => (
  <div>
    <h1>{ gameState }</h1>
  </div>
);

const mapStateToProps = (state) => {
  const game = state.entities.find(e => e.tags.type === 'game')

  if(!game) {
    return { gameState: "unknown" }
  }

  return {
    gameState: game.tags.state
  }
}

GameState = connect(mapStateToProps)(GameState);

export default GameState;
