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
    try {
      this.dogs = await getJson('dog');
    } catch (e) {
      this.handleError(e);
    }
  },
  methods: {
    clear() {
      this.id = null;
      this.breed = '';
      this.name = '';
    },
    createDog() {
      const {breed, name} = this;
      this.handlePromise(postJson('dog', {breed, name}), async res => {
        const dog = await res.json();
        this.dogs.push(dog);
      });
    },
    deleteDog(id) {
      this.handlePromise(deleteResource(`dog/${id}`), () => {
        const index = this.dogs.findIndex(dog => dog.id === id);
        this.$delete(this.dogs, index);
      });
    },
    handleError(e) {
      const msg = e instanceof Error ? e.message : e;
      alert(msg);
    },
    async handlePromise(promise, callback) {
      try {
        const res = await promise;
        if (res.ok) {
          callback(res);
          this.clear();
        } else {
          this.handleError(await res.text());
        }
      } catch (e) {
        this.handleError(e);
      }
    },
    selectDog(dog) {
      this.breed = dog.breed;
      this.id = dog.id;
      this.name = dog.name;
    },
    updateDog(dog) {
      const {breed, id, name} = this;
      dog.id = id;
      dog.breed = breed;
      dog.name = name;
      this.handlePromise(putJson(`dog/${id}`, dog), () => {
        const index = this.dogs.findIndex(dog => dog.id === id);
        this.$set(this.dogs, index, dog);
      });
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
