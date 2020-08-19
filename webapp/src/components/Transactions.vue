<template>
  <table class="table is-striped is-fullwidth">
    <thead>
      <tr>
        <th>Date</th>
        <th>Entity</th>
        <th>Amount</th>
        <th>Note</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="transaction in transactions" v-bind:key="transaction._id">
        <td>{{ transaction.date }}</td>
        <td>{{ transaction.entity }}</td>
        <td>{{ transaction.amount }}</td>
        <td>{{ transaction.note }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component
export default class Transactions extends Vue {
  data() {
    return {
      transactions: [
        // TODO: consider modifying this to be a prop that's passed in
        {
          date: '08/04/2020',
          entity: 'Amazin',
          amount: 200,
          note: 'beep boop',
        },
        {
          date: '08/04/2020',
          entity: 'Amazin',
          amount: 200,
          note: 'beep boop',
        },
      ],
    };
  }
}
</script>
