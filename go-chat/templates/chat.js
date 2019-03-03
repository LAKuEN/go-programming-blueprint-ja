let socket = null;
const form = document.getElementById("#chatbox"), 
  msgBox = document.querySelector("#chatbox > textarea"),
  messages = document.getElementById("#messages");

form.addEventListener("submit", (ev) => {
  console.log(msgBox.text);
  if (!msgBox.text) {
    return false;
  }
  if (!socket) {
    alert("Error: There is no connection with WebSocket");
    return false;
  }
  socket.send(msgBox.text);
  msgBox.text = "";
  return false;
});

if (!window["websocket"]) {
  const msg = "Error: You are using a browser that WebSocket is not available";
  alert(msg);
  throw new Error(msg);
}
let socket = new WebSocket("ws://localhost:8080/room");
socket.onclose = () => {
  alert("Disconnected");
}
socket.onmessage = (ev) => {
  newMsg = document.createElement("li");
  newMsg.textContent = ev.data;

  messages.appendChild(newMsg);
}
