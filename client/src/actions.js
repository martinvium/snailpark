export const receiveMessageAction = (dispatch, event) => {
  const message = JSON.parse(event.data)
  console.log(message)
  dispatch({ type: message.t, ...message.m })
}

export const repeatingPingAction = (connection) => {
  return (dispatch, getState) => {
    const { playerId } = getState()

    connection.send(JSON.stringify({
      action: 'ping',
      playerId: playerId
    }));

    setTimeout(() => dispatch(repeatingPingAction), 45 * 1000);
  }
}

export const startAction = (connection) => {
  return (dispatch, getState) => {
    const { playerId } = getState()

    const message = JSON.stringify({
      action: 'start',
      playerId: playerId
    })

    connection.send(message);
  }
}

export const clickCardAction = (entityId) => {
  return (dispatch, getState) => {
    const { connection, playerId } = getState()

    const message = JSON.stringify({
      action: 'playCard',
      playerId: playerId, // TODO why do we need to include this?
      card: entityId
    })

    connection.send(message);
  }
}
