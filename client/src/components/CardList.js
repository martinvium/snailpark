import React from 'react';
import { connect } from 'react-redux';
import Card from './Card'

let CardList = ({ entities }) => {
  return (
    <div>
      {entities.map(e => (
        <Card key={e.id} entity={e} />
      ))}
    </div>
  )
}


const mapStateToProps = (state) => {
  return {
    entities: state.entities
  }
}

CardList = connect(mapStateToProps)(CardList);

export default CardList;
