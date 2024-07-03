# snake-ai

Snake AI is a **competitive self-education game** specifically crafted for **programming learning**.

Building AI for your Snake you can play with your friends and improve your algorithms and coding skills.

<div align="center">
  <img class="logo" src="https://github.com/andreamper220/snakeai/assets/55195085/7d51c629-83d3-42b7-82eb-e25d3cf902ba" width="400px" alt="Snake-AI"/>
</div>

<div align="center">
  [![Go Report Card](https://goreportcard.com/badge/github.com/andreamper220/snakeai)](https://goreportcard.com/report/github.com/andreamper220/snakeai)&nbsp;
</div>

## What is it and how it works?

Snake-AI helps you and your friends to exercise in coding and/or algorythms.

You can:
- **Register** with your email and password
- **Login** to start playing
- **Create** your lobby as **HOST**:
  - Up to **10** players simultaneously
  - From **5x5** to **30x30** field size
- **Connect** to any existing lobby as **CLIENT** according to your game skills

### How does Match Making work?

As a **HOST**:
1. You **create** a new lobby
2. Your data is starting to process by separate goroutine:
   1. If you are a single player - you **create a party** with yourself and **start to play**
   2. Else - you are waiting for another players, increasing your skill delta **twice** every **3 sec**

As a **CLIENT**:
1. You **request to connect** to any existing lobby
2. Your data is starting to process by several separate goroutines:
   1. If there is any party with `party.members.skill.avg = [ you.skill - you.delta; you.skill + you.delta ]`, then - you connect to that party
   2. Else - you are waiting for some party, increasing your skill delta **twice** every **3 sec**

## What can I do while playing?

You can:
- **Write AI** with following commands:
  - **Left** - turn head left
  - **Right** - turn head right
  - **Move** - continue moving
- **Try your AI**: unleash your snake AI power on the game field!
- **Play with friends**: who is the best with coding skills? Let's know it!

_**NB**_: Equally **1 snake** is possible for **1 player**!

_**NNB**_: Up to **10 players** are possible for **1 game session**!

## Installation

- `make build`
- `make up`
- ENJOY :)
