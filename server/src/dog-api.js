const express = require('express');
const api = express.Router();

const PgConnection = require('postgresql-easy');
const dbConfig = {database: 'survey'};
const pg = new PgConnection(dbConfig);

function handleError(res, action, error) {
  const msg = error ? `failed to ${action}: ${error.message}` : action;
  console.error(msg);
  res.status(500).send(msg);
}

api.post('/', async (req, res) => {
  const {breed, name} = req.body;

  try {
    const dog = {breed, name};
    const id = await pg.insert('dog', {breed, name});
    if (id) {
      dog.id = id;
      res.send(dog);
    } else {
      handleError(res, 'create dog');
    }
  } catch (e) {
    handleError(res, 'create dog', e);
  }
});

api.get('/', async (req, res) => {
  try {
    const dogs = await pg.query('select * from dog');
    res.send(dogs);
  } catch (e) {
    handleError(res, 'get dogs', e);
  }
});

api.put('/:id', async (req, res) => {
  const {id} = req.params;
  const {breed, name} = req.body;

  try {
    await pg.updateById('dog', id, {breed, name});
    res.sendStatus(200);
  } catch (e) {
    handleError(res, 'update dog', e);
  }
});

api.delete('/:id', async (req, res) => {
  const {id} = req.params;
  try {
    await pg.deleteById('dog', id);
    res.sendStatus(200);
  } catch (e) {
    handleError(res, `delete dog ${id}`, e);
  }
});

module.exports = api;
