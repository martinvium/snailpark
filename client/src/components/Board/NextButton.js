import React from 'react'
import { connect } from 'react-redux'
import { nextAction } from '../../actions'

let NextButton = ({ onClick }) => {
  let classNames = ['btn']
  if (false) {
    classNames.push('btn-disabled')
  }

  return (
    <input type="button" id="end-turn" value="Next" onClick={onClick} className={classNames.join(' ')}/>
  )
}

const mapStateToProps = (state, ownProps) => {
  const { gameState } = state

  let text = 'End turn'
  if(['blockers', 'blockTarget'].includes(gameState)) {
    text = 'Begin combat'
  }

  return { text }
}

const mapDispatchToProps = (dispatch) => {
  return {
    onClick: () => (
      dispatch(nextAction())
    )
  }
}

NextButton = connect(mapStateToProps, mapDispatchToProps)(NextButton);

export default NextButton;
