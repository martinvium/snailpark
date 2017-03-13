$(document).ready(function() {
  var card_proto = $('#card-prototype'),
    stack = [],
    playerId = 'player';

  var players = {
    "player": { "id": "player", "hand": [], "board": [], "currentMana": 0, "maxMana": 0 },
    "ai": { "id": "ai", "hand": [], "board": [], "currentMana": 0, "maxMana": 0 },
  };

  var ws = new WebSocket(getWebSocketUrl('entry'));

  ws.onopen = function(event) {
    ws.send(
      JSON.stringify({
        "playerId": playerId,
        "action": "start"
      })
    );
  }

  ws.onmessage = function(event) {
    var msg = JSON.parse(event.data);

    players = msg.players;

    clearBoard();
    renderBoard();
    renderHand();
    renderMana();
  }

  $('#end-turn').click(function() {
    console.log('End turn');
    ws.send(JSON.stringify({ "playerId": playerId, "action": "end_turn" }));
  });

  function boardEl(playerId) {
    return $('#' + playerId + ' .board');
  }

  function handEl(playerId) {
    return $('#' + playerId + ' .hand');
  }

  function manaEl(playerId) {
    return $('#' + playerId + ' .mana');
  }

  function playCard(id) {
    console.log('Playing card: ' + id);

    ws.send(
      JSON.stringify({
        "playerId": playerId,
        "action": "play_card",
        cards: [{ "id": id }]
      })
    );
  }

  function clearBoard() {
    $.each(players, function(playerId, _) {
      boardEl(playerId).empty();
      handEl(playerId).empty();
    });
  }

  function renderBoard() {
    $.each(players, function(_, player) {
      $.each(player["board"], function(index, card) {
        boardEl(player["id"]).append(
          renderCard(card)
        );
      });
    });
  }

  function renderHand() {
    $.each(players, function(_, player) {
      $.each(player["hand"], function(index, card) {
        handEl(player["id"]).append(
          renderCard(card, function() {
            playCard($(this).attr('data-id'));
          })
        );
      });
    });
  }

  function renderCard(value, callback) {
    console.log('Rendering: ' + value.title);
    var card = card_proto.clone();
    card.attr('data-id', value.id);
    card.addClass('orange').show();
    card.click(callback);
    $('.title', card).text(value['title']);
    $('.cost', card).text(value['cost']);
    $('.type', card).text(value['type']);
    $('.description', card).text(value['description']);
    return card;
  }

  function renderMana() {
    $.each(players, function(_, player) {
      $('.current', manaEl(player["id"])).text(player["currentMana"]);
      $('.max', manaEl(player["id"])).text(player["maxMana"]);
    });
  }

  function getWebSocketUrl(s) {
    var l = window.location;
    return ((l.protocol === "https:") ? "wss://" : "ws://") + l.host + l.pathname + s;
  }
});

