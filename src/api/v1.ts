import express from 'express';
const router = express.Router();
import mysql from 'mysql2';
import passport from 'passport';
const { EduPgae } = require('edupage-api-forweb');
const CronJob = require('cron').CronJob;
import puppeteer from 'puppeteer';
import { randomAsciiString, simpleStringify, replaceCircularFast } from '../utils/misc';
import { EP2User } from '../utils/user';
import { exit } from 'process';

if (process.env.DATABASE_URL == undefined) {
  console.log('Please define the URL of the MySQL database in the .env file')
  exit(1)
}
const connection = mysql.createConnection(process.env.DATABASE_URL);
(async function () {
  const browser = await puppeteer.launch({
    executablePath: '/usr/bin/chromium-browser',
    args: ['--no-sandbox', '--headless']
  });
})()

const users = new Map();

router.post('/auth', async (req, res) => {
  if(req.body.token) {
    const u = await EP2User.fromToken(req.body.token, connection);
    if(u === null || typeof u === 'boolean') {
      res.status(401).json({
        error: 'Invalid token',
        valid: false,
        token: req.body.token
      })
    } else {
      users.set(u.token, u)
      res.status(200).json({
        valid: true,
        token: req.body.token
      })
    }
  } else if(req.body.email && req.body.password) {
    const u = await EP2User.fromEmailPassword(req.body.email, req.body.password, connection);
    if(u === null || typeof u === 'boolean') {
      res.status(401).json({
        error: 'Not found',
        valid: false,
        token: null
      })
    } else {
      users.set(u.token, u)
      res.status(200).json({
        valid: true,
        token: u.token
      })
    }
  } else {
    res.status(400).json({
      error: 'Please supply token for validation or email/password for authentication'
    })
  }
});

router.post('/register', async (req, res) => {
  const u = await EP2User.createFromEmailPassword(req.body.email, req.body.password, connection);
  res.status(200).send(u.token)
})

module.exports = router;