<template>
  <section class="hero">
    <div class="hero-body">
      <div class="container">
        <h1 class="title is-1">
          Budgeter
        </h1>
      </div>
    </div>
  </section>
  <!-- <nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
      <div class="navbar-item content ml-3 mt-3">
        <h1>Budgeter</h1>
      </div>

      <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false">
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </a>
    </div>
  </nav> -->
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Navbar extends Vue {
  @Prop() private msg!: string;
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.hero {
  background-color: #95b8d1;
  font-family: 'PT Serif', serif;
}
</style>
