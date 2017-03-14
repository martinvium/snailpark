$(document).ready(function() {
  var card_proto = $('#card-prototype'),
    stack = [],
    playerId = 'player';

  var players = {
    "player": { "id": "player", "hand": [], "board": [], "currentMana": 0, "maxMana": 0 },
    "ai": { "id": "ai", "hand": [], "board": [], "currentMana": 0, "maxMana": 0 },
  };

  var gameId = urlParam('guid');
  if(!gameId) {
    gameId = guid();
  }

  var ws;

  openWebsocket();

  $('#end-turn').click(function() {
    console.log('End turn');
    ws.send(JSON.stringify({ "playerId": playerId, "action": "end_turn" }));
  });

  $('#reconnect').click(function() {
    console.log('Reconnect');
    ws.close();
  });

  function openWebsocket() {
    ws = new WebSocket(getWebSocketUrl('game/connect?gameId=' + gameId));

    ws.onopen = function(event) {
      $('#messages').text('').removeClass().hide();

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
      renderHealth();

      if(msg.state == "finished") {
        if(players["player"]["health"] <= 0) {
          $('#messages').text('You lost! :(').addClass('red').show();
        } else {
          $('#messages').text('You won! :)').addClass('green').show();
        }
      }
    }

    ws.onerror = function(event) {
      console.log("ERROR", event);
    }

    ws.onclose = function(event) {
      $('#messages').text('Disconnected! Reconnecting...').addClass('yellow').show();
      setTimeout(openWebsocket, 5000);
    }
  }

  function boardEl(playerId) {
    return $('#' + playerId + ' .board');
  }

  function handEl(playerId) {
    return $('#' + playerId + ' .hand');
  }

  function manaEl(playerId) {
    return $('#' + playerId + ' .mana');
  }

  function healthEl(playerId) {
    return $('#' + playerId + ' .health');
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

  function renderHealth() {
    $.each(players, function(_, player) {
      $('.current', healthEl(player["id"])).text(player["health"]);
    });
  }

  function getWebSocketUrl(s) {
    var l = window.location;
    return ((l.protocol === "https:") ? "wss://" : "ws://") + l.host + l.pathname + s;
  }

  function guid() {
    function s4() {
      return Math.floor((1 + Math.random()) * 0x10000)
        .toString(16)
        .substring(1);
    }
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
      s4() + '-' + s4() + s4() + s4();
  }

  function urlParam(sParam) {
    var sPageURL = window.location.search.substring(1);
    var sURLVariables = sPageURL.split('&');
    for (var i = 0; i < sURLVariables.length; i++) {
      var sParameterName = sURLVariables[i].split('=');
      if (sParameterName[0] == sParam) {
        return sParameterName[1];
      }
    }
  }
});

