<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
  </div>
</template>

<script>
import axios from "axios";
export default {
  data: function() {
    return {
      notes: []
    };
  },
  created: function() {
    this.getNotes();
  },
  methods: {
    getNotes: function() {
      axios
        .get("localhost:9003/api/status")
        .then(response => (this.notes = response.data))
        .catch(error => console.log(error.response.data));
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
