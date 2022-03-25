<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <h1>{{ dateTime }}</h1>
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

    onMounted(async () => {
      const response = await fetch('http://localhost:13011/get-time', {
        method: 'GET'
      })

      const json = await response.json()
      const date = new Date(json.time)
      dateTime.value = date.toLocaleString()
      console.log("get-time: " + dateTime.value)

      eventSource.addEventListener("time", function (event) {
        const json = JSON.parse(event.data)
        const date = new Date(json.time)
        dateTime.value = date.toLocaleString()
        console.log("time-event: " + dateTime.value)
      })
    })

    return {
      dateTime,
      eventSource
    }
  },
}
</script>

<style scoped>
h3 {
  margin: 40px 0 0;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>