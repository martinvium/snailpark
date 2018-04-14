const initialState = {
  gameState: "test",
  playerId: "player",
  entities: [],
  players: {}
}

export const appReducer = (state = initialState, action) => {
  switch(action.type) {
    case 'FULL_STATE':
      return fullStateReducer(state, action)
    case 'CHANGE_ATTR':
      return { ...state, entities: entitiesReducer(state.entities, action) }
    case 'CHANGE_TAG':
      return { ...state, entities: entitiesReducer(state.entities, action) }
    case 'REVEAL_ENTITY':
      return revealEntityReducer(state, action)
    default:
      return state
  }
}

// updateState();
export const fullStateReducer = (state, action) => {
  const { entities, players } = action

  return {
    ...state,
    player: players["player"],
    enemy: players["ai"],
    entities
  }
}

export const entitiesReducer = (state, action) => {
  switch(action.type) {
    case 'CHANGE_ATTR':
    case 'CHANGE_TAG':
      return state.map(e => entityReducer(e, action))
    default:
      return state
  }
}

export const entityReducer = (state, action) => {
  const { Key, Value } = action

  if(state.id !== action.entityId) {
    return state
  }

  switch(action.type) {
    case 'CHANGE_ATTR':
      let { attributes } = state
      attributes[Key] = Value
      return { ...state, attributes }
    case 'CHANGE_TAG':
      if(state.id !== action.entityId) {
        return state
      }

      let { tags } = state
      tags[Key] = Value
      return { ...state, tags }
    default:
      return state
  }
}

export const revealEntityReducer = (state, action) => {
  // handleRevealEntity(packet.m);
  // updateState();
  return state
}
