import {createApp} from 'vue'
import App from './App.vue'

const app = createApp(App)
const eventSource = new EventSource("http://localhost:13011/time-event", {withCredentials: false});

app.provide('eventSource', eventSource)
app.mount('#app')