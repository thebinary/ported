<template>
  <div>
    <div v-for="(value, key) in data" :key="key">
      <span class="primary-text has-text-weight-bold"> {{ key }}</span>
      <span class="">: {{ value[0] }}</span>
    </div>
  </div>
</template>

<script>
export default {
  name: "Headers",
  props: {
    data: {
      type: Object,
    },
  },
};
</script>

<style></style>
