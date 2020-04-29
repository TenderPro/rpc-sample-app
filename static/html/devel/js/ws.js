
var wsFactory = new Map();  //  websocket store
var divFactory = new Map(); //  div store

function setDiv(tag) {
  divFactory.set(tag, document.getElementById(tag));
}

function pageLoaded(useFiles) {
  setDiv("time");
  setDiv("echo");
  setDiv("log");
  setDiv("err");
  var f = document.getElementById("time.form");
  sendForm(f);
}

function sendForm(form) {
  var op = form.getAttribute('op');
  var arg = form.elements["arg"].value;
  console.log("form_op:", op, "arg:", arg);
  switch (op) {
    case "time":
      openStream(op, arg, true, function(data){
        var tm = data.result.ts;
        var tms = new Date(tm).toLocaleString();
        divFactory.get("time").innerHTML = tms;
      });
      break;
    case "echo":
      var div = divFactory.get("echo");
      div.innerHTML = '';
      openStream("sample/ping/list", arg, false, function(data){
        div.innerHTML += "<br>" + data.result.counter+": "+data.result.Value;
      });
      break;
    default:
      console.log('Unknown op: '+op);
  }
  return false;
}

function openStream(op, arg, restart, cb) {
  console.log("STREAM> op:", op, "arg:", arg);

  var loc = window.location, url;
  var url = ((loc.protocol === "https:")?"ws":"w") + "s://" + loc.host + "/v1/"+op;
  if (arg !== undefined) {
    url += "/"+arg
  }
  if (wsFactory.has(op)) {
    // Close old connection
    wsFactory.get(op).close();
  }
  var log = divFactory.get("log");
  var err = divFactory.get("err");
  var ws = new WebSocket(url);
  wsFactory.set(op, ws);
  err.innerHTML = '';
  ws.onmessage = function(msg) {
    console.log('WS>>>', op, arg, msg.data)
    var data = JSON.parse(msg.data);
    if (data.result === undefined) {
      err.innerHTML += data.error;
    } else {
      cb(data);
    }
  };
  ws.onclose = function(event) {
    console.debug('Code: ' + event.code + ' reason: ' + event.reason);
    if (event.wasClean) {
      console.debug('Connection closed clean');
    } else {
      log.innerHTML=op+': Connection closed';
      console.debug(op, 'Connection closed');
      if (restart) {
        // Try to reconnect in 5 seconds
        setTimeout(function(){openStream(op, arg, cb)}, 5000);
      }
    }
  };
  ws.onopen = function(){
    log.innerHTML='';
  };
  ws.onerror = function(evt) {
    console.log("WS ERROR: " + evt.data)
  };
}

