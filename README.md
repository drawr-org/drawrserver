# drawr / core-server

[![Build Status](https://jenkins.etsag.de/buildStatus/icon?job=drawr-core-server-linux)](https://jenkins.etsag.de/job/drawr-core-server-linux/)

The backend to the drawr service

# API

## GET: `/session/new`
Request a new session

## GET: `/session/[session-id]`
Get session information  

## ws:// `/session/[session-id]/ws`
Connect to a session websocket
