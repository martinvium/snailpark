import React from 'react';

const Card = (props) => {
  const { id, tags, anonymous } = props

  var value = tags.attackTarget;
  const attacking = typeof value !== 'undefined' && value !== '';

  let classNames = [ 'card' ]
  classNames.push(tags.color)

  if(anonymous) {
    classNames.push('card-back')
  }

  if(attacking) {
    classNames.push('attacking')
  }

  return (
    <div className={classNames.join(' ')} id={id}>
      <CardBody { ...props } />
    </div>
  )
}

export default Card

const artUrl = type => (
  require(`../../images/${type}.png`)
)

const CardBody = ({ id, anonymous, tags, attributes, damage, onMouseOver, onMouseOut, onClick }) => {
  if(anonymous) {
    return null
  }

  return (
    <div onMouseOver={e => onMouseOver(e, id)} onMouseOut={e => onMouseOut(e, id)} onClick={e => onClick({ id })}>
      <div className="header">
        <div className="title pull-left">{ tags.title }</div>
        <div className="cost pull-right">{ attributes.cost }</div>
        <div className="cf"></div>
      </div>
      <div className="art" style={{backgroundImage: `url(${artUrl(tags.type)})`}}></div>
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
  if(typeof damage === 'undefined') {
    return null
  }

  return (
    <div className="baloon">-{ damage }</div>
  )
}
