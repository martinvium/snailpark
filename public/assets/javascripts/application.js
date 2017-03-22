$(document).ready(function() {
  var card_proto = $('#card-prototype'),
    playerId = 'player',
    players = {},
    oldPlayers = null,
    state,
    msg;

  var gameId = urlParam('guid');
  if(!gameId) {
    gameId = guid();
  }

  var ws;

  var pingTime = 45 * 1000;
  var ping;

  openWebsocket();

  $('#end-turn').click(function() {
    ws.send(JSON.stringify({ "playerId": playerId, "action": "endTurn" }));
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
      msg = JSON.parse(event.data);

      state = msg.state;

      if(oldPlayers == null) {
        oldPlayers = msg.players;
      } else {
        oldPlayers = players;
      }
      players = msg.players;

      updateStateHelp(msg);
      clearBoard();
      renderBoard();
      renderHand();
      renderPlayers();
      renderEngagementArrows();

      if(msg.state == "finished") {
        if(getAvatar(players["player"])["currentToughness"] <= 0) {
          $('#messages').text('You lost! :(').addClass('red').show();
        } else {
          $('#messages').text('You won! :)').addClass('green').show();
        }
      }

      $.each(players, function(key, player) {
        var old = getAvatar(oldPlayers[key])["currentToughness"],
            now = getAvatar(player)["currentToughness"];

        if(now < old) {
          $('.current', healthEl(player["id"])).effect('highlight', { color: 'red' }, 3000);
        } else if (now > old) {
          $('.current', healthEl(player["id"])).effect('highlight', { color: 'cyan' }, 3000);
        }
      });

      updateNextButton(msg);

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

  $(window).mousemove(function(e) {
    if(state == 'blockTarget') {
      renderPointerArrow(e);
    }
  });

  function renderPointerArrow(e) {
    var pos = getCardPosition(msg.currentBlocker['id']);

    var path = $('.arrow svg path.pointer');
    if(path.length === 0) {
      path = clonePathPrototype().addClass('pointer');
      $('.arrow svg').append(path);
    }

    path.attr('d', getArrowPathD(
      pos['x'],
      pos['y'],
      e.clientX,
      e.clientY
    ));
  }

  function renderEngagementArrows() {
    $('.arrow svg path.engagement, .arrow svg path.pointer').remove();

    for(var i in msg.engagements) {
      var eng = msg.engagements[i];
      if(!eng['blocker'] || !eng['attacker']) {
        continue;
      }

      var startPos = getCardPosition(eng['blocker']['id']);
      var targetPos = getCardPosition(eng['attacker']['id']);

      path = clonePathPrototype().addClass('engagement');

      path.attr('d', getArrowPathD(
        startPos['x'],
        startPos['y'],
        targetPos['x'],
        targetPos['y']
      ));

      $('.arrow svg').append(path);
    }
  }

  function clonePathPrototype() {
    return $('.arrow svg path.prototype').clone().removeClass('prototype');
  }

  function getCardPosition(id) {
    var card = $('[data-id="' + id + '"]');
    var offset = card.offset();
    var x = offset.left + card.width() / 2;
    var y = offset.top + (card.height() / 4 * 1);
    return { 'x': x, 'y': y };
  }

  // TODO: Curve
  // var p1 = 'C326,212';
  // var p2 = '900,900';
  function getArrowPathD(sx, sy, tx, ty) {
    var start = 'M' + sx + ',' + sy;
    var target = tx + ',' + ty;
    var points = [start, target];
    return points.join(' ');
  }

  function updateNextButton(msg) {
    if(msg.currentPlayerId == playerId) {
      $('#end-turn').removeClass('btn-disabled');
    } else {
      $('#end-turn').addClass('btn-disabled');
    }
  }

  function getAvatar(player) {
    for(var i in player["board"]) {
      if(player["board"][i]["type"] === "avatar") {
        return player["board"][i];
      }
    }

    return { "currentToughness": 0 };
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

  function clickCard(id, action) {
    console.log('Playing card: ' + id);

    ws.send(
      JSON.stringify({
        "playerId": playerId,
        "action": action,
        "card": id
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
          renderCard(card, function() {
            clickCard($(this).attr('data-id'), 'target');
          })
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
              clickCard($(this).attr('data-id'), 'playCard');
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
    card.addClass(value['color']).show();
    card.click(callback);
    $('.title', card).text(value['title']);
    $('.cost', card).text(value['cost']);
    $('.type', card).text(value['type']);
    $('.description', card).text(value['description']);
    $('img', card).attr('src', '/assets/images/' + value['type'] + '.jpg');
    if(value['toughness'] > 0) {
      $('.power-toughness .power', card).text(value['power']);
      $('.power-toughness .toughness', card).text(value['currentToughness']);
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
      var avatar = getAvatar(player);
      $('.current', healthEl(player["id"])).text(avatar["currentToughness"]);
    });
  }

  function updateStateHelp(msg) {
    var el = $('#state-help');

    if(msg.currentPlayerId != playerId) {
      el.text('Your opponent is thinking...')
      return
    }

    if(state == "main") {
      el.text('Play a card, or attack with your creatures by clicking on them...');
    } else if(state == "targeting") {
      el.text('Pick a target...');
    } else if(state == "blockers" || state == "blockTarget") {
      el.text('Opponent is attacking! Assign blockers and end turn when you are ready');
    }
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

