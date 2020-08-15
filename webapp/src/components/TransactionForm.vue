<template>
  <!-- TODO: Add a field for the date, or some kind of extra settings -->
  <div class="columns is-1 is-variable">
    <div class="column content is-2">
      <h1>Date</h1>
      <input class="input" type="date" placeholder="Text input" />
    </div>
    <div class="column content">
      <h1>Entity</h1>
      <input class="input" type="text" placeholder="Text input" />
    </div>
    <div class="column content is-2">
      <h1>Amount</h1>
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
export default class TransactionForm extends Vue {}
</script>
