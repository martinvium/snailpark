import React from 'react';
import { connect } from 'react-redux';
import Card from './Card'

let CardList = ({ id, entities }) => {
  return (
    <div>
      {entities.map(e => (
        <Card key={e.id} entity={e} />
      ))}
    </div>
  )
}


const mapStateToProps = (state, ownProps) => {
  const entities = state.entities.filter(entity => (
    entity.playerId === ownProps.id && entity.tags.location === ownProps.location
  ))

  return {
    id: ownProps.id,
    entities 
  }
}

CardList = connect(mapStateToProps)(CardList);

export default CardList;
