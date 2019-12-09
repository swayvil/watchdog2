import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import 'bootstrap'
import 'bootstrap/dist/css/bootstrap.min.css'
import './assets/styles/main.css'
import LazyLoadDirective from "./directives/LazyLoadDirective";

Vue.config.productionTip = false;
axios.defaults.baseURL = 'http://localhost:9999';

Vue.directive("lazyload", LazyLoadDirective);

new Vue({
  render: h => h(App),
}).$mount('#app')
