import React from 'react';
import { connect } from 'react-redux';
import CardList from '../CardList'

let PlayerTerritory = () => {
  return (
    <div>
      <CardList/>
    </div>
  )
}

PlayerTerritory = connect()(PlayerTerritory);

export default PlayerTerritory;
