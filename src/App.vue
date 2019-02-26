<template>
  <Dogs />
</template>

<script>
import Dogs from './components/Dogs.vue';

export default {
  name: 'App',
  components: {Dogs}
};
</script>

<style>
body {
  background-color: cornflowerblue;
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  text-align: center;
}

button:disabled {
  color: gray;
}

label {
  font-size: 18px;
  color: white;
}
</style>
