<template>
  <div id="logger">
    <div>
      <b v-if="connected" class="active">Active</b>
      <b v-else class="inactive">Inactive</b>
    </div>
    <br />
    <div class="log" v-for="item in logs" v-bind:key="item">{{ item }}</div>
  </div>
</template>

<script>
const maxRecords = 50;
export default {
  name: "Logger",
  data: function() {
    return {
      connected: false,
      logs: [],
      stream: null,
    };
  },
  created: function() {
    this.stream = new WebSocket(
      "ws://" + window.location.host + "/porter/stream"
    );

    this.stream.onopen = function() {
      this.connected = true;
      console.log("Connected");
    }.bind(this);

    this.stream.onmessage = function(event) {
      let data = JSON.parse(event.data);
      //this.logs.unshift(data);
      this.logs.push(data);
      if (this.logs.length > maxRecords) {
        this.logs.shift();
        //this.logs.pop();
      }
    }.bind(this);

    this.stream.onclose = function() {
      console.log("Disconnected");
      this.connected = false;
    }.bind(this);
  },
};
</script>

<style>
.active {
  color: green;
}
.inactive {
  color: #ff0000d6;
}
</style>
