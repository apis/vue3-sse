<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <div id="wrapper">
      <div/>
      <div/>
      <div>
        <h1 class="time" :class="actualColor">{{ dateTime }}</h1>
      </div>
      <div/>
      <div/>
    </div>
  </div>
</template>

<script>
import {inject, onMounted, ref} from 'vue'

export default {
  name: 'HelloWorld',
  props: {
    msg: String
  },
  setup() {
    const eventSource = inject('eventSource')
    const dateTime = ref("")
    const primaryColor = ref(true)
    const actualColor = ref("")

    onMounted(async () => {
      const response = await fetch('/api/get-time', {
        method: 'GET'
      })

      const json = await response.json()
      const date = new Date(json.time)
      dateTime.value = date.toLocaleString()

      function changeColor() {
        if (!primaryColor.value) {
          actualColor.value = "blue"
        } else {
          actualColor.value = "green"
        }
        primaryColor.value = !primaryColor.value
      }

      changeColor()
      console.log("get-time: " + dateTime.value)



      eventSource.addEventListener("time", function (event) {
        const json = JSON.parse(event.data)
        const date = new Date(json.time)
        dateTime.value = date.toLocaleString()
        changeColor()
        console.log("time-event: " + dateTime.value)
      })
    })

    return {
      actualColor,
      primaryColor,
      dateTime,
      eventSource
    }
  },
}
</script>

<style scoped>
#wrapper {
  display: flex;
  /*height: 200px;*/
}

#wrapper > div {
  flex-grow: 1;
}

.time {
  padding: 5px;
  transition: all 1s;
}

.blue {
  color: white;
  background: blue;
  opacity: 0.7;
}

.green {
  color: white;
  background: green;
  opacity: 0.7;
}
</style>
