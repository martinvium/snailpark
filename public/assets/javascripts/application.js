$(document).ready(function() {
  var card_proto = $('#card-prototype'),
    ui_board = $('#me .board').html(''),
    ui_hand = $('#me .hand').html(''),
    hand = [],
    board = [],
    stack = [],
    life = 100,
    mana = 100;

  var ws = new WebSocket("ws://localhost:8080/entry");

  ws.onopen = function(event) {
    ws.send(JSON.stringify({ "action": "start" }));
  }

  ws.onmessage = function(event) {
    var msg = JSON.parse(event.data);

    // msg also has state of the FSM and possibly priority e.g. you or me
    // msg also has a stack_resolved: true param, which means the user action
    // resolves and the ui can empty the stack?
    switch(msg.action) {
      case "update":
        // could represent start_turn, or a transition between states etc.
        break;
      case "add_to_hand": // draw
        console.log("Add to hand!");
        for(var i in msg.cards) {
          hand.push(msg.cards[i]);
        }
        break;
      case "put_on_stack":
        for(var i in msg.cards) {
          for(var y in hand) {
            if(hand[y].id == msg.cards[i].id) {
              hand.splice(y, 1);
              break;
            }
          }

          stack.push(msg.cards[i]);
        }
        break;
      case "empty_stack":
        stack = [];
        break;
      case "add_to_board":
        for(var i in msg.cards) {
          board.push(msg.cards[i]);
        }
        break;
      case "hit_player":
        break;
      case "remove_from_board":
        msg.cards // array of cards to be removed from board
        break;
      case "action_not_valid":
        // cancel the user action in the ui
        break;
      case "pick_card": // a card you used triggered an event where you must select a card from a list of choices
        msg.cards // array of cards to present for selection
        break;
    }

    clearBoard();
    renderBoard();
    renderHand();
  }

  $('#end-turn').click(function() {
    console.log('End turn');
    ws.send(JSON.stringify({ "action": "end_turn" }));
  });

  function playCard(id) {
    console.log('Playing card: ' + id);
    ws.send(JSON.stringify({ "action": "play_card", cards: [{ "id": id }] }));
  }

  function clearBoard() {
    ui_board.empty();
    ui_hand.empty();
  }

  function renderBoard() {
    $.each(board, function(index, value) {
      var card = renderCard(value);
      ui_board.append(card);
    });
  }

  function renderHand() {
    $.each(hand, function(index, value) {
      var card = renderCard(value, function() {
        playCard($(this).attr('data-id'));
      });

      ui_hand.append(card);
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

