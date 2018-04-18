import React from 'react';
import { connect } from 'react-redux';
import CardDetails from './CardDetails'
import PlayerTerritory from './PlayerTerritory'

let Board = ({ targeting }) => {
  let classNames = []

  if(targeting) {
    classNames.push('targeting')
  }

  return (
    <div id="board" className={classNames}>
      <CardDetails/>

      <PlayerTerritory id={`ai`}/>
      <PlayerTerritory id={`player`}/>
    </div>
  )

  // return (
  //   <div id="board" className={classNames}>
  //     <ModalDialog/>
  //     <NextButton/>
  //     <CardDetails/>

  //     <PlayerTerritory type={`opponent`}/>
  //     <CombatTerritory/>
  //     <PlayerTerritory type={`yours`}/>
  //   </div>
  // )
}


const mapStateToProps = (state) => {
  return {
    targeting: false
  }
}

Board = connect(mapStateToProps)(Board);

export default Board;
