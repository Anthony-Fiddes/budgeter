<template>
  <section class="px-4 pt-2">
    <div class="box">
      <TransactionForm />
      <Transactions />
    </div>
  </section>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import TransactionForm from './TransactionForm.vue';
import Transactions from './Transactions.vue';

@Component({
  components: { TransactionForm, Transactions },
})
export default class MainView extends Vue {}
</script>
