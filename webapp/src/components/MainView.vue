<template>
  <section class="px-4 pt-2">
    <div class="box">
      <TransactionForm />
    </div>
  </section>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import TransactionForm from './TransactionForm.vue';
import ItemList from './ItemList.vue';

@Component({
  components: { TransactionForm },
})
export default class MainView extends Vue {}
</script>
