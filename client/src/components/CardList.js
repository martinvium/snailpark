import React from 'react';
import { connect } from 'react-redux';
import { clickCardAction } from '../actions'
import Card from './Card'

let CardList = ({ id, entities, onMouseOver, onMouseOut, onClick }) => {
  return (
    <div>
      {entities.map(e => (
        <Card key={e.id} { ...e } onMouseOver={onMouseOver} onMouseOut={onMouseOut} onClick={onClick}  />
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

const mapDispatchToProps = (dispatch) => {
  return {
    onMouseOver: (e, entityId) => {
      dispatch({ type: "UPDATE_CARD_DETAILS", id: entityId })
    },
    onMouseOut: (e) => {
      dispatch({ type: "UPDATE_CARD_DETAILS", id: null })
    },
    onClick: (e, entityId) => {
      dispatch(clickCardAction(entityId))
    }
  }
}

CardList = connect(mapStateToProps, mapDispatchToProps)(CardList);

export default CardList;
