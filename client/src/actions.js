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
