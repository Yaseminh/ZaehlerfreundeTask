# challenge.zaehlerfreunde.com

## Prerequisites

1 - Golang (v 1.22.0 or later) on the Backend

2 - Node.js version 8.9 or above (v10+ recommended) to run Vuejs on the Frontend

## Setup

To start the Go service run

```
cd cost-calculator
go run .
```

For the Vue application execute the following commands

```
cd zaehlerfreunde-challenge-project
npm install
npm run dev
```

## Instructions

The goal of the challenge is to calculate yesterday's energy cost based on the consumption measured by an electricity meter and the energy price from the day-ahead EPEX sport market.

An electricity meter measures energy consumption in kWh. A meter reading is the measured energy consumption at a specific time. (https://en.wikipedia.org/wiki/Electricity_meter)

In contrast to private households, commercial consumers can buy electricity directly at the energy market. We want to calculate the cost for a commercial consumer that buys its electricity from the EPEX spot Day-Ahead market. It specifies a fixed energy price for each hour of the day.

#### API Endpoint

We want to create a REST endpoint that accepts a list of 25 meter readings, one for each hour of the day + one for midnight of the next day, and returns the energy cost.

In order to get the market prices we will use the Awattar REST API GET endpoint:

- https://api.awattar.de/v1/marketdata?start=[START-TIMESTAMP]&end=[END-TIMESTAMP]

- Timestamps should be in milliseconds (i.e 1724857124000). This website https://www.epochconverter.com/ can be helpful.

#### Web App

The Vue app should show the meter readings in a table and send them to the API endpoint to calculate their energy cost.

In App.vue the function `generateYesterdaysReadings` generates random meter readings for each hour of yesterday's date. You will need to add:

- A table which displays the meter readings
- A function to call the API endpoint that's called when the _Calculate Cost_ button is clicked.
