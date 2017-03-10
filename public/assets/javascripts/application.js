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

    players[msg.playerId]["currentMana"] = msg["currentMana"];
    players[msg.playerId]["maxMana"] = msg["maxMana"];
    console.log(players[msg.playerId]);

    // msg also has state of the FSM and possibly priority e.g. you or me
    // msg also has a stack_resolved: true param, which means the user action
    // resolves and the ui can empty the stack?
    switch(msg.action) {
      case "add_to_hand": // draw
        players[msg.playerId]["hand"] = $.merge(
          players[msg.playerId]["hand"],
          msg.cards
        );
        break;
      case "put_on_stack":
        // Remove cards from hand (really this verbose?)
        $.each(msg.cards, function(_, card) {
          $.each(players[msg.playerId]["hand"], function(index, cardInHand) {
            if(cardInHand.id == card.id) {
              players[msg.playerId]["hand"].splice(index, 1);
              return false;
            }
          });
        });

        stack = $.merge(stack, msg.cards);
        break;
      case "empty_stack":
        stack = [];
        break;
      case "add_to_board":
        players[msg.playerId].board = $.merge(
          players[msg.playerId].board,
          msg.cards
        );
        break;
    }

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
    $.each(players, function(_, client) {
      $.each(client["board"], function(index, card) {
        boardEl(client["id"]).append(
          renderCard(card)
        );
      });
    });
  }

  function renderHand() {
    $.each(players, function(_, client) {
      $.each(client["hand"], function(index, card) {
        handEl(client["id"]).append(
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
    $('.type', card).text(value['type']);
    $('.description', card).text(value['description']);
    return card;
  }

  function renderMana() {
    $.each(players, function(_, client) {
      $('.current', manaEl(client["id"])).text(client["currentMana"]);
      $('.max', manaEl(client["id"])).text(client["maxMana"]);
    });
  }

  function getWebSocketUrl(s) {
    var l = window.location;
    return ((l.protocol === "https:") ? "wss://" : "ws://") + l.host + l.pathname + s;
  }
});

