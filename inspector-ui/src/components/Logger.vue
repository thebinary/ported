<template>
  <div id="logger">
    <div>
      <b v-if="connected" class="active">Active</b>
      <b v-else class="inactive">Inactive</b>
    </div>
    <br />
    <b-collapse
      class="card"
      animation="slide"
      v-for="item in logs"
      :key="item.id"
      :open="item.isOpen == true"
      @open="item.isOpen = true"
    >
      <div slot="trigger" slot-scope="props" class="card-header" role="button">
        <div class="card-header-title">
          <div class="r-ip">{{ item.ip }}</div>
          <div class="r-method">
            <span class="method" :class="item.m">{{ item.m + " " }}</span>
            <span>{{ item.p }}</span>
          </div>
          <div class="r-status">{{ item.s }}</div>
          <div class="r-cLength">{{ item.rpl }} B</div>
          <div class="r-rtt">{{ item.rtt }}</div>
        </div>
        <a class="card-header-icon">
          <b-icon
            :icon="props.open ? 'arrow-down' : 'arrow-up'"
            class="is-size-7 primary-text"
          >
          </b-icon>
        </a>
      </div>
      <div class="card-content">
        <p class="is-uppercase has-text-weight-bold is-size-6">
          Request Headers
        </p>
        <Headers :data="item.rh" />
        <br />
        <p class="is-uppercase has-text-weight-bold is-size-6">
          Response Headers
        </p>
        <Headers :data="item.rph" />
      </div>
    </b-collapse>
  </div>
</template>

<script>
import Headers from "./Headers.vue";
const maxRecords = 50;
const webSocketURL = "ws://" + window.location.host + "/porter/stream";
export default {
  name: "Logger",
  components: {
    Headers,
  },
  data: function() {
    return {
      connected: false,
      isOpen: -1,
      logs: [],
      stream: null,
    };
  },
  created: function() {
    this.stream = new WebSocket(webSocketURL);

    this.stream.onopen = function() {
      this.connected = true;
      console.log("Connected");
    }.bind(this);

    this.stream.onmessage = function(event) {
      let d = JSON.parse(event.data);
      let data = JSON.parse(d);
      data.id = data.t + "-" + data.rtt;
      data.isOpen = false;
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
.method {
  font-weight: bold;
}
.GET {
  color: green;
}
.POST {
  color: orange;
}
.PATCH {
  color: darkorange;
}
.DELETE {
  color: red;
}
.request-time {
  color: darkgray;
}
.r-method {
  flex-grow: 8;
}
.r-ip {
  flex-grow: 1;
}
.r-status {
  flex-grow: 1;
}
.r-cLength {
  flex-grow: 1;
}
</style>
