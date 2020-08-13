<template>
  <div id="app">
    <Navbar msg="Welcome to Your Vue.js + TypeScript App" />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Navbar from './components/Navbar.vue';

@Component({
  components: {
    Navbar,
  },
})
export default class App extends Vue {}
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=PT+Serif|Open+Sans&display=swap');
html {
  background-color: #424c55;
}
</style>
