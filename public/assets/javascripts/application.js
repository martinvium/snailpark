$(document).ready(function() {
  var card_proto = $('#card-prototype'),
    stack = [],
    playerClientId = 'player';

  var clients = {
    "player": { "id": "player", "hand": [], "board": [] },
    "ai": { "id": "ai", "hand": [], "board": [] },
  };

  var ws = new WebSocket("ws://localhost:8080/entry");

  ws.onopen = function(event) {
    ws.send(
      JSON.stringify({
        "clientId": playerClientId,
        "action": "start"
      })
    );
  }

  ws.onmessage = function(event) {
    var msg = JSON.parse(event.data);

    // msg also has state of the FSM and possibly priority e.g. you or me
    // msg also has a stack_resolved: true param, which means the user action
    // resolves and the ui can empty the stack?
    switch(msg.action) {
      case "add_to_hand": // draw
        clients[msg.clientId]["hand"] = $.merge(
          clients[msg.clientId]["hand"],
          msg.cards
        );
        break;
      case "put_on_stack":
        // Remove cards from hand (really this verbose?)
        $.each(msg.cards, function(_, card) {
          $.each(clients[msg.clientId]["hand"], function(index, cardInHand) {
            if(cardInHand.id == card.id) {
              clients[msg.clientId]["hand"].splice(index, 1);
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
        clients[msg.clientId].board = $.merge(
          clients[msg.clientId].board,
          msg.cards
        );
        break;
    }

    clearBoard();
    renderBoard();
    renderHand();
  }

  $('#end-turn').click(function() {
    console.log('End turn');
    ws.send(JSON.stringify({ "clientId": playerClientId, "action": "end_turn" }));
  });

  function boardEl(playerClientId) {
    return $('#' + playerClientId + ' .board');
  }

  function handEl(playerClientId) {
    return $('#' + playerClientId + ' .hand');
  }

  function playCard(id) {
    console.log('Playing card: ' + id);

    ws.send(
      JSON.stringify({
        "clientId": playerClientId,
        "action": "play_card",
        cards: [{ "id": id }]
      })
    );
  }

  function clearBoard() {
    $.each(clients, function(clientId, _) {
      boardEl(clientId).empty();
      handEl(clientId).empty();
    });
  }

  function renderBoard() {
    $.each(clients, function(_, client) {
      $.each(client["board"], function(index, card) {
        boardEl(client["id"]).append(
          renderCard(card)
        );
      });
    });
  }

  function renderHand() {
    $.each(clients, function(_, client) {
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
});

