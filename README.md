# drawrserver

[![Build Status](https://jenkins.etsag.de/buildStatus/icon?job=drawr-core-server-linux)](https://jenkins.etsag.de/job/drawr-core-server-linux/)

The backend of the drawr service

# API documentation

<details>
<summary>/session</summary>

- `/new` :: **GET** :: requests a new session

</details>

<details>
<summary>/session/:sessionID</summary>

- `/` :: **GET** :: returns session information
- `/` :: **POST** :: updates session information
- `/` :: **DELETE** :: delete a session from the database
- `/ws` :: **GET** :: websocket of the session
- `/leave` :: **GET** :: disconnect from websocket *deprecated*

</details>

<details>
<summary>/stats</summary>

- `/` :: **GET** :: statistics report for the server
- `/db` :: **GET** :: statistics report for the database

</details>

<details>
<summary>/version</summary>

- `/` :: **GET** :: returns API version (also used for connection testing)

</details>
