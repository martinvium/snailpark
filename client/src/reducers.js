const initialState = {
  gameState: "test",
  playerId: "player",
  entities: [],
  players: {},
  cardDetails: null
}

export const appReducer = (state = initialState, action) => {
  switch(action.type) {
    case 'FULL_STATE':
      return fullStateReducer(state, action)
    case 'CHANGE_ATTR':
    case 'CHANGE_TAG':
    case 'REVEAL_ENTITY':
      return { ...state, entities: entitiesReducer(state.entities, action) }
    case 'UPDATE_CARD_DETAILS':
      return { ...state, cardDetails: action.id }
    case 'UPDATE_CONNECTION':
      return { ...state, connection: action.connection }
    default:
      return state
  }
}

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
    case 'REVEAL_ENTITY':
      if(state.some(e => e.id === action.entityId)) {
        return state.map(e => entityReducer(e, { ...action, type: 'UPDATE_ENTITY' }))
      } else {
        return [...state, entityReducer(undefined, {...action, type: 'ADD_ENTITY'})]
      }
    default:
      return state
  }
}

export const entityReducer = (state, action) => {
  switch(action.type) {
    case 'CHANGE_ATTR':
      if(state.id !== action.entityId) {
        return state
      }

      let { attributes } = state
      attributes[action.Key] = action.Value
      return { ...state, attributes }
    case 'CHANGE_TAG':
      if(state.id !== action.entityId) {
        return state
      }

      let { tags } = state
      tags[action.Key] = action.Value
      return { ...state, tags }
    case 'UPDATE_ENTITY':
      if(state.id !== action.entityId) {
        return state
      }

      return {
        ...state,
        ...action.entity
      }
    case 'ADD_ENTITY':
      return {
        ...state,
        ...action.entity
      }
    default:
      return state
  }
}
