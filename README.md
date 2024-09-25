# snake-ai

Snake AI is a **competitive self-education game** specifically crafted for **programming learning**.

Building AI for your Snake you can play with your friends and improve your algorithms and coding skills.

<div align="center">
  <img class="logo" src="https://github.com/andreamper220/snakeai/assets/55195085/7d51c629-83d3-42b7-82eb-e25d3cf902ba" width="400px" alt="Snake-AI"/>
</div>

<div align="center">

[![build](https://github.com/andreamper220/snakeai/actions/workflows/ci.yml/badge.svg)](https://github.com/andreamper220/snakeai/actions/workflows/ci.yml)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/andreamper220/snakeai)](https://goreportcard.com/report/github.com/andreamper220/snakeai)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/andreamper220/snakeai/badge.svg?branch=master)](https://coveralls.io/github/andreamper220/snakeai?branch=main)

</div>

## What is it and how it works?

Snake-AI helps you and your friends to exercise in coding and/or algorythms.

You can:
- **Register** with your email and password
- **Login** to start playing
- **Create** your lobby as **HOST**:
  - Up to **10** players simultaneously
  - From **5x5** to **30x30** field size
  - (optional) Create your **CUSTOM** map with our map editor!
- **Connect** to any existing lobby as **CLIENT** according to your game skills

### How does Match Making work?

Here is the full UML diagram to explain Match Making:

![RLLHRZ~1](https://github.com/andreamper220/snakeai/assets/55195085/f807595d-7c9b-4d5b-bdaf-831112d04b11)

## What can I do while playing?

You can:
- **Write AI** with following:
  - Commands:
    - **Left** - turn head left
    - **Right** - turn head right
    - **Move** - continue moving
  - Condition Operators:
    - **If**
    - **Else If** (you can chain them multiple times)
    - **Else**
  - Conditions (**>, <, >=, <=, ==, !=**):
    - **Distance to Edge / Wall**
    - **Distance to Food**
    - **Distance to another Snake**
- **Try your AI**: unleash your snake AI power on the game field!
- **Play with friends**: who is the best with coding skills? Let's know it!

_**NB**_: Equally **1 snake** is possible for **1 player**!

_**NNB**_: Up to **10 players** are possible for **1 game session**!

## Installation

- `make build`
- `make up`
- ENJOY :)
