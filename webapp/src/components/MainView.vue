<template>
  <section class="px-4 pt-2">
    <div class="box">
      <ItemForm />
    </div>
  </section>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import ItemForm from './ItemForm.vue';

@Component({
  components: { ItemForm },
})
export default class MainView extends Vue {}
</script>
