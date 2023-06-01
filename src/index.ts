import express from 'express';
import passport from 'passport';
const Strategy = require('passport-http-bearer').Strategy;
import mysql from 'mysql2';
import { exit } from 'process';
require('dotenv').config();

if (process.env.DATABASE_URL == undefined) {
  console.log('Please define the URL of the MySQL database in the .env file')
  exit(1)
}
const connection = mysql.createConnection(process.env.DATABASE_URL)
const app = express();
require('express-ws')(app);

passport.use(new Strategy(
  async function (token: String, cb: Function) {
    connection.query('SELECT * FROM users WHERE token = ?', [token], function (err: mysql.QueryError | null, results: Array<mysql.RowDataPacket>, fields: Array<mysql.FieldPacket>) {
      if (err) {
        console.log(err);
        return cb(err);
      }
      if (results.length > 0) {
        return cb(null, results[0]);
      } else {
        return cb(null, false);
      }
    });
  }
));

app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(require('morgan')('combined'));
app.use(require('cors')());

app.use(require('./api/v1'))

app.get('/', (_, res) => {
  res.send('Public API server for EduPage2')
});

app.listen(8080, () => {
  console.log('Server started on port 8080')
})