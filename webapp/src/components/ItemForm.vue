<template>
  <!-- TODO: Add a field for the date, or some kind of extra settings -->
  <div class="columns is-1 is-variable item-form">
    <div class="column content">
      <h1>Item</h1>
      <input class="input" type="text" placeholder="Text input" />
    </div>
    <div class="column content">
      <h1>Vendor</h1>
      <input class="input" type="text" placeholder="Text input" />
    </div>
    <div class="column content is-1">
      <h1>Cost</h1>
      <input class="input" type="text" placeholder="Text input" />
    </div>
    <div class="column content">
      <h1>Note</h1>
      <input class="input" type="text" placeholder="Text input" />
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component
export default class ItemForm extends Vue {}
</script>
