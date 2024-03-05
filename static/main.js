let ws = null;

(function () {
  const url = "ws://localhost:8080/ws";
  const timeoutSecs = 1000;
  let tryCount = 0;
  const canvasDiv = document.getElementById("canvas_div");
  const canvas = document.getElementById("canvas");
  const ctx = canvas.getContext("2d");
  //   const size = 10 * 10;
  let lastWidth = 0;
  const cellArray = [];

  console.log("canvasDiv", canvasDiv);
  console.log("canvas", canvas);

  const drawHtml = (data) => {
    if (data.width !== lastWidth) {
      console.log(data.width, data.size);

      lastWidth = data.width;

      for (let index = 0; index < cellArray.length; index++) {
        cellArray[index].remove();
      }
      cellArray.length = 0;

      for (let idx = 0; idx < data.board.length; idx++) {
        const cell = document.createElement("i");

        cell.style.width = `${100 / data.width}%`;

        canvasDiv.appendChild(cell);
        cellArray.push(cell);
      }
    }

    for (let idx = 0; idx < data.board.length; idx++) {
      if (data.board[idx] === 1) {
        cellArray[idx].style.backgroundColor = "#333";
        // cellArray[idx].setAttribute("class", "alive");
      } else {
        cellArray[idx].style.backgroundColor = "#ccc";
        // cellArray[idx].setAttribute("class", "dead");
      }
    }
  };

  const drawCanvas = (data) => {
    if (data.width !== lastWidth) {
      console.log(data.width, data.size);

      lastWidth = data.width;

      canvas.setAttribute("width", data.width * 10);
      canvas.setAttribute("height", data.width * 10);

      // canvas.setAttribute("width", data.width * 1);
      // canvas.setAttribute("height", data.width * 1);
    }

    for (let idx = 0; idx < data.board.length; idx++) {
      if (data.board[idx] === 1) {
        ctx.fillStyle = "#333";
      } else {
        ctx.fillStyle = "#ccc";
      }

      const x = idx % data.width;
      const y = Math.floor(idx / data.width);

      ctx.fillRect(x * 10, y * 10, 10, 10);
      // ctx.fillRect(x * 1, y * 1, 1, 1);
    }
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

      if (tryCount <= 5) {
        tryCount += 1;
      }
    });

    ws.addEventListener("message", (e) => {
      try {
        drawCanvas(JSON.parse(e.data));
        // drawHtml(JSON.parse(e.data));
      } catch (err) {
        console.error(err);
        console.log("message", e);
      }
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
