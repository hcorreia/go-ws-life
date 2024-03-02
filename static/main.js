let ws = null;

(function () {
  const url = "ws://localhost:8080/ws";
  const timeoutSecs = 1000;
  let tryCount = 0;
  const canvas = document.getElementById("canvas");
  //   const size = 10 * 10;

  console.log("canvas", canvas);

  const draw = (data) => {
    let result = "";

    for (let idx = 0; idx < data.length; idx++) {
      if (data[idx] === 1) {
        result += '<i class="alive"></i>';
      } else {
        result += '<i class="dead"></i>';
      }
    }

    canvas.innerHTML = result;
  };

  const wsConnect = () => {
    console.log("connect", tryCount);

    if (ws && ws.readyState !== WebSocket.CLOSED) {
      return;
    }

    ws = new WebSocket(url);

    ws.addEventListener("open", (e) => {
      console.log("open", e);

      //   timeoutSecs = 1000;
      tryCount = 0;

      ws.send("PING");
    });

    ws.addEventListener("close", (e) => {
      console.log("close", e);

      console.log("timeout", timeoutSecs * tryCount * 2);
      setTimeout(wsConnect, timeoutSecs * tryCount * 2);

      if (tryCount <= 10) {
        tryCount += 1;
      }
    });

    ws.addEventListener("message", (e) => {
      //   console.log("message", e);

      draw(JSON.parse(e.data));
    });

    ws.addEventListener("error", (e) => {
      console.log("error", e);

      //   setTimeout(wsConnect, timeoutSecs);
    });

    window.ws = ws;
  };

  wsConnect();

  //   window.addEventListener("load", wsConnect);
  //   document.addEventListener("load", wsConnect);
})();
