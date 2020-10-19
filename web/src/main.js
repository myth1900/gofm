import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import axios from 'axios'
import store from "@/store";
Vue.use(store)

Vue.config.productionTip = false
if (process.env.VUE_APP_MODE==="development"){
  axios.defaults.baseURL = process.env.VUE_APP_API_URL
}
Vue.prototype.$axios = axios
new Vue({
  el: '#app',
  store,
  vuetify,
  render: h => h(App)
})
