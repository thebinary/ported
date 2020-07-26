<template>
  <div id="app">
    <Logger />
  </div>
</template>

<script>
import Logger from "./components/Logger.vue";

export default {
  name: "App",
  components: {
    Logger,
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-left: 20px;
  margin-top: 20px;
}
</style>
