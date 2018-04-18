import React from 'react';
import CardList from '../CardList'

let PlayerTerritory = ({ id }) => {
  let classNames = [id, "cl"]

  return (
    <div>
      <div className="battlefield-wrapper">
        <div className="battlefield">
          <div className={classNames.join(' ')}>
            <CardList id={id} location="board"/>
          </div>
        </div>
      </div>
      <div className={classNames.join(' ')}>
        <ul className="hand cl">
          <CardList id={id} location="hand"/>
        </ul>
      </div>
    </div>
  )
}

export default PlayerTerritory;
