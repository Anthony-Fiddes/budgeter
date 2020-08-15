<template>
  <div id="app">
    <Navbar />
    <MainView />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Navbar from './components/Navbar.vue';
import MainView from './components/MainView.vue';

@Component({
  components: {
    Navbar,
    MainView,
  },
})
export default class App extends Vue {}
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=PT+Serif|Open+Sans&display=swap');
html {
  /* the !important is required for Chrome */
  background-color: #424c55 !important;
}
</style>
