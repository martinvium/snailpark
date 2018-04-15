import React from 'react';
import { connect } from 'react-redux';

export default ({ entity, onClickCard }) => {
  var value = entity.tags.attackTarget;
  const attacking = typeof value != 'undefined' && value != '';

  let classNames = [ 'card' ]
  classNames.push(entity.tags.color)

  if(entity.anonymous) {
    classNames.push('card-back')
  }

  if(attacking) {
    classNames.push('attacking')
  }

  return (
    <div className={classNames.join(' ')} id={entity.id}>
      <CardBody {...entity } />
    </div>
  )
}

const CardBody = ({ id, anonymous, tags, attributes, damage, onClickCard }) => {
  if(anonymous) {
    return null
  }

  return (
    <div onClick={e => onClickCard({ id })}>
      <div className="header">
        <div className="title pull-left">{ tags.title }</div>
        <div className="cost pull-right">{ attributes.cost }</div>
        <div className="cf"></div>
      </div>
      <div className="art" style={{backgroundImage: `url('/assets/images/${tags.type}.png')`}}></div>
      <div className="description">{ tags.description }</div>
      <div className="power-toughness">
        <PowerLabel {...attributes} />
        <span className="type center">{ tags.type }</span>
        <ToughnessLabel {...attributes} />
      </div>
      <DamageBaloon damage={damage} />
    </div>
  )
}

const PowerLabel = ({ power, toughness }) => {
  // console.log(power, toughness)
  if(typeof toughness === 'undefined' || typeof power === 'undefined') {
    return null
  }

  return (
    <span className="power pull-left">{ power }</span>
  )
}

const ToughnessLabel = ({ toughness }) => {
  if(typeof toughness === 'undefined') {
    return null
  }

  return (
    <span className="toughness pull-right">{ toughness }</span>
  )
}

const DamageBaloon = ({ damage }) => {
  console.log(damage)
  if(typeof damage === 'undefined') {
    return null
  }

  return (
    <div className="baloon">-{ damage }</div>
  )
}
