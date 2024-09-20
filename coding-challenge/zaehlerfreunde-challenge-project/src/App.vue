<template>
  <div id="app">
    <div class="btns">
      <button @click="generateYesterdaysReadings()">Generate readings</button>
      <button @click="calculateCost">Calculate cost</button>
    </div>

    <!-- Table displaying meter readings -->
    <table v-if="meterReadings.length" class="meter-table">
      <thead>
      <tr>
        <th>Timestamp</th>
        <th>kWh</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(reading, index) in meterReadings" :key="index">
        <td>{{ reading.timestamp.toLocaleString() }}</td>
        <td>{{ reading.kwh.toFixed(2) }}</td>
      </tr>
      </tbody>
    </table>

    <!-- Display the calculated cost -->
    <div>
      Cost: €{{ cost.toLocaleString('en', { maximumFractionDigits: 2, minimumFractionDigits: 2 }) }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, type Ref } from 'vue'
import axios from 'axios'

// Interface for meter readings
interface MeterReading {
  timestamp: Date
  kwh: number
}

// Reactive data for meter readings and calculated cost
const meterReadings: Ref<MeterReading[]> = ref([])
const cost: Ref<number> = ref(0)

// Generates random meter readings for each hour of yesterday
function generateYesterdaysReadings() {
  var kwh = Math.random() * 1000000
  const readings = []

  for (var hour = 0; hour <= 24; hour++) {
    const date = new Date()
    date.setDate(date.getDate() - 1)
    date.setHours(hour, 0, 0, 0)

    kwh += Math.random() * 100

    readings.push({
      timestamp: date,
      kwh: kwh
    })
  }

  meterReadings.value = readings
}

// Function to calculate energy cost by calling the API endpoint
async function calculateCost() {
  if (!meterReadings.value.length) return

  // Prepare the readings for the API request (timestamps in ms and kWh values)
  const formattedReadings = meterReadings.value.map(reading => ({
    timestamp: reading.timestamp.getTime(),
    value: reading.kwh
  }))

  try {
    // Explicitly specify the backend URL (http://localhost:8080)
    const response = await axios.post('http://localhost:8080/energy_cost', {
      readings: formattedReadings
    })

    // Set the returned total cost from the API response
    cost.value = response.data.total_cost
  } catch (error) {
    console.error('Error calculating cost:', error)
    alert('Error calculating cost, please try again.')
  }
}
</script>

<style lang="css" scoped>
#app {
  padding: 30px;
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.btns {
  display: flex;
  gap: 8px;
}

.meter-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

.meter-table th, .meter-table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: center;
}

.meter-table th {
  background-color: #f4f4f4;
}
</style>
