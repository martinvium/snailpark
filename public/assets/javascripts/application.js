$(document).ready(function() {
  var card_proto = $('#card-prototype'),
    stack = [],
    playerId = 'player';

  var players = {},
    oldPlayers = null;

  var gameId = urlParam('guid');
  if(!gameId) {
    gameId = guid();
  }

  var ws;

  var pingTime = 45 * 1000;
  var ping;

  openWebsocket();

  $('#end-turn').click(function() {
    ws.send(JSON.stringify({ "playerId": playerId, "action": "end_turn" }));
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

      if(oldPlayers == null) {
        oldPlayers = msg.players;
      } else {
        oldPlayers = players;
      }
      players = msg.players;

      clearBoard();
      renderBoard();
      renderHand();
      renderPlayers();

      if(msg.state == "finished") {
        if(players["player"]["health"] <= 0) {
          $('#messages').text('You lost! :(').addClass('red').show();
        } else {
          $('#messages').text('You won! :)').addClass('green').show();
        }
      }

      $.each(players, function(key, player) {
        if(player["health"] < oldPlayers[key]["health"]) {
          $('.current', healthEl(player["id"])).effect('highlight', { color: 'red' }, 3000);
        } else if (player['health'] > oldPlayers[key]['health']) {
          $('.current', healthEl(player["id"])).effect('highlight', { color: 'cyan' }, 3000);
        }
      });

      if(msg.currentPlayerId == playerId) {
        $('#end-turn').removeClass('btn-disabled');
      } else {
        $('#end-turn').addClass('btn-disabled');
      }

      clearTimeout(ping);
      ping = setTimeout(pingSocket, pingTime);
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

  function pingSocket(id) {
    ws.send(JSON.stringify({
      "playerId": playerId,
      "action": "ping"
    }));

    ping = setTimeout(pingSocket, pingTime);
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
        if(card.type != "creature") {
          console.log("skipping non creature");
          return;
        }

        boardEl(player["id"]).append(
          renderCard(card)
        );
      });
    });
  }

  function renderHand() {
    $.each(players, function(_, player) {
      if(player["id"] == playerId) {
        $.each(player["hand"], function(index, card) {
          handEl(player["id"]).append(
            renderCard(card, function() {
              playCard($(this).attr('data-id'));
            })
          );
        });
      } else {
        for(var i = 0; i < player["handSize"]; i++) {
          handEl(player["id"]).append(
            renderCardBack()
          );
        }
      }
    });
  }

  function renderPlayers() {
    renderMana();
    renderHealth();
  }

  function renderCardBack(value, callback) {
    var card = card_proto.clone();
    card.text('').addClass('card-back').show();
    return card;
  }

  function renderCard(value, callback) {
    var card = card_proto.clone();
    card.attr('data-id', value.id);
    card.addClass('orange').show();
    card.click(callback);
    $('.title', card).text(value['title']);
    $('.cost', card).text(value['cost']);
    $('.type', card).text(value['type']);
    $('.description', card).text(value['description']);
    $('img', card).attr('src', '/assets/images/' + value['type'] + '.jpg');
    if(value['type'] == 'creature') {
      $('.power-toughness .power', card).text(value['power']);
      $('.power-toughness .toughness', card).text(value['toughness']);
    } else {
      $('.power-toughness', card).hide();
    }
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

