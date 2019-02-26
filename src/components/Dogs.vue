<template>
  <div class="dogs">
    <h1>Dogs</h1>
    <div>
      <label>Breed</label>
      <input type="text" v-model="breed" />
    </div>
    <div>
      <label>Name</label>
      <input type="text" v-model="name" />
    </div>
    <div>
      <button @click="createDog" :disabled="!breed || !name">Create</button>
      <button @click="updateDog" :disabled="!id">Update</button>
    </div>
    <table>
      <tr>
        <th>Name</th>
        <th>Breed</th>
        <th>Delete</th>
      </tr>
      <tr v-for="dog in dogs" :key="dog.id" @click="selectDog(dog)">
        <td>{{ dog.name }}</td>
        <td>{{ dog.breed }}</td>
        <td><button @click="deleteDog(dog.id)">X</button></td>
      </tr>
    </table>
  </div>
</template>

<script>
import {deleteResource, getJson, postJson, putJson} from '../fetch-util';
export default {
  name: 'Dogs',
  data() {
    return {
      breed: '',
      dogs: [],
      id: null,
      name: ''
    };
  },
  async mounted() {
    this.dogs = await getJson('dog');
  },
  methods: {
    async createDog() {
      const {breed, name} = this;
      try {
        const res = await postJson('dog', {breed, name});
        if (res.ok) {
          const dog = await res.json();
          this.dogs.push(dog);
          this.id = null;
          this.breed = '';
          this.name = '';
        } else {
          this.handleError(res.text);
        }
      } catch (e) {
        this.handleError(e);
      }
    },
    async deleteDog(id) {
      try {
        await deleteResource(`dog/${id}`);
        const index = this.dogs.findIndex(dog => dog.id === id);
        this.$delete(this.dogs, index);
      } catch (e) {
        this.handleError(e);
      }
    },
    async getDogs() {
      try {
        this.dogs = await getJson('dog');
      } catch (e) {
        this.handleError(e);
      }
    },
    handleError(e) {
      const msg = e instanceof Error ? e.message : e;
      alert(msg);
    },
    selectDog(dog) {
      this.breed = dog.breed;
      this.id = dog.id;
      this.name = dog.name;
    },
    async updateDog(dog) {
      const {breed, id, name} = this;
      dog.id = id;
      dog.breed = breed;
      dog.name = name;
      try {
        const res = await putJson(`dog/${id}`, dog);
        if (res.ok) {
          const index = this.dogs.findIndex(dog => dog.id === id);
          this.$set(this.dogs, index, dog);
          this.id = null;
          this.breed = '';
          this.name = '';
        } else {
          this.handleError(res.text);
        }
      } catch (e) {
        this.handleError(e);
      }
    }
  }
};
</script>

<style scoped>
button {
  background-color: white;
  border-radius: 4px;
  color: black;
  margin: 0 5px 10px 0;
  padding: 4px;
}

div {
  margin-top: 10px;
}

.dogs {
  font-size: 18px;
}

h1 {
  color: white;
}

input {
  border-radius: 4px;
  padding: 4px;
}

label {
  margin-right: 10px;
}

table {
  color: white;
  margin: 0 auto;
}

td,
th {
  border: solid white 1px;
  padding: 10px;
  text-align: center;
}
</style>
