<template>
  <div id="app">
    <nav id="navbar" class="navbar">
      <div class="navbar-brand">
        <div class="navbar-item has-text-light">
          <b-icon id="app-title-icon" icon="binoculars"></b-icon>
          <p id="app-title" class="is-size-6 has-text-weight-bold">
            Ported Inspector
          </p>
        </div>
      </div>
    </nav>
    <div id="app-container">
      <Logger />
    </div>
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
  //margin: 20px;
}
#navbar {
  background-color: #236993;
}
#app-container {
  margin: 20px;
}
#app-title-icon {
  margin-right: 0.3em;
}
.primary-text {
  color: #236993;
}
</style>
