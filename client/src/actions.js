export const receiveMessageAction = (dispatch, event) => {
  const message = JSON.parse(event.data)
  console.log(message)
  dispatch({ type: message.t, ...message.m })
}

export const repeatingPingAction = (socket) => {
  return (dispatch, getState) => {
    const { playerId } = getState()

    socket.send(JSON.stringify({
      action: 'ping',
      playerId: playerId
    }));

    setTimeout(() => dispatch(repeatingPingAction), 45 * 1000);
  }
}

export const startAction = (socket) => {
  return (dispatch, getState) => {
    const { playerId } = getState()

    const message = JSON.stringify({
      action: 'start',
      playerId: playerId
    })

    socket.send(message);
  }
}
  // socket.on('FULL_STATE', data => fullStateAction(dispatch, data))
  // socket.on('CHANGE_ATTR', data => changeAttrAction(dispatch, data))
  // socket.on('CHANGE_TAG', data => changeTagAction(dispatch, data))
  // socket.on('REVEAL_ENTITY', data => revealEntityAction(dispatch, data))
// }

// export const fullStateAction = (dispatch, data) => {
  // dispatch({ type: "FULL_STATE", ...data })
// }

// export const changeAttrAction = (dispatch, data) => {
  // dispatch({ type: "CHANGE_ATTR", ...data })
// }

// export const changeTagAction = (dispatch, data) => {
  // dispatch({ type: "CHANGE_TAG", ...data })
// }

// export const revealEntityAction = (dispatch, data) => {
  // dispatch({ type: "REVEAL_ENTITY", ...data })
// }
