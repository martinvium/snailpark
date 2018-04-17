import React from 'react';
import { connect } from 'react-redux';
import Card from './Card'

let CardList = ({ id, entities }) => {
  let classNames = [id, "cl"]

  return (
    <div className={classNames.join(' ')}>
      {entities.map(e => (
        <Card key={e.id} entity={e} />
      ))}
    </div>
  )
}


const mapStateToProps = (state, ownProps) => {
  return {
    id: ownProps.id,
    entities: state.entities.filter(entity => entity.playerId === ownProps.id)
  }
}

CardList = connect(mapStateToProps)(CardList);

export default CardList;
