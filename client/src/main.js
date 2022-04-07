import {createApp} from 'vue'
import App from './App.vue'

const app = createApp(App)
const eventSource = new EventSource("/api/time-event", {withCredentials: false});

app.provide('eventSource', eventSource)
app.mount('#app')
