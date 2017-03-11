# drawrserver

[![Build Status](https://jenkins.etsag.de/buildStatus/icon?job=drawr-core-server-linux)](https://jenkins.etsag.de/job/drawr-core-server-linux/)

The backend of the drawr service

# API

* **GET** `/session/new` - request a new session

* **GET** `/session/<session-id>` - retrieve session information  

* **ws://** `/session/<session-id>/ws` - connect to a session websocket
