const express = require('express');
const api = express.Router();

const PgConnection = require('postgresql-easy');
const dbConfig = {database: 'survey', user: 'postgres'};
const pg = new PgConnection(dbConfig);

async function tryIt(res, action, fn) {
  try {
    let result = await fn();
    const sendResult = action.startsWith('get') || action.startsWith('create');
    res.status(200).send(sendResult ? result : null);
  } catch (e) {
    const msg = 'failed to ' + action + (e ? `: ${e.message}` : '');
    console.error(msg);
    res.status(500).send(msg);
  }
}

api.post('/', (req, res) => {
  const {breed, name} = req.body;
  const dog = {breed, name};
  tryIt(res, 'create dog', async () => {
    dog.id = await pg.insert('dog', dog);
    if (!dog.id) throw new Error('failed to insert dog');
    return dog;
  });
});

api.get('/', (req, res) =>
  tryIt(res, 'get dogs', () => pg.query('select * from dog'))
);

api.put('/:id', (req, res) => {
  const {id} = req.params;
  const {breed, name} = req.body;
  tryIt(res, 'update dog', () => pg.updateById('dog', id, {breed, name}));
});

api.delete('/:id', (req, res) => {
  const {id} = req.params;
  tryIt(res, 'delete dog', () => pg.deleteById('dog', id));
});

module.exports = api;
