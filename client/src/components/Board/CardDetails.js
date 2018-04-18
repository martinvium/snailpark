import React from 'react';
import { connect } from 'react-redux';
import Card from '../Card'

let CardDetails = ({ entity }) => {
  if (!entity) {
    return null
  }

  return (
    <div className="card-details">
      <Card { ...entity } />
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    entity: state.entities.find(e => state.cardDetails === e.id) 
  }
}

CardDetails = connect(mapStateToProps)(CardDetails);

export default CardDetails;
