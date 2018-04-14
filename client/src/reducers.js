const initialState = {
  playerId: "player"
}

export const appReducer = (state = initialState, action) => {
  switch(action.type) {
    case 'FULL_STATE':
      return fullStateReducer(state, action)
    case 'CHANGE_ATTR':
      return changeAttrReducer(state, action)
    case 'CHANGE_TAG':
      return changeTagReducer(state, action)
    case 'REVEAL_ENTITY':
      return revealEntityReducer(state, action)
    default:
      return state
  }
}

export const fullStateReducer = (state, action) => {
  // handleFullState(packet.m);
  // updateState();
  return state
}

export const changeAttrReducer = (state, action) => {
  // handleChangeAttrTag(packet.m, 'attributes');
  // updateState();
  return state
}

export const changeTagReducer = (state, action) => {
  // handleChangeAttrTag(packet.m, 'tags');
  // updateState();
  return state
}

export const revealEntityReducer = (state, action) => {
  // handleRevealEntity(packet.m);
  // updateState();
  return state
}
