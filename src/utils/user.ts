import mysql, { RowDataPacket } from 'mysql2';
import puppeteer, { Browser, Page } from 'puppeteer';
import { randomAsciiString } from './misc';

export class EP2User {
  email: string = ""
  password: string = ""
  icanteen_email: string | null
  icanteen_password: string | null
  token: string = ""
  page: Page | null
  connection: mysql.Connection
  last_active: Date | null

  constructor(email: string, password: string, icanteen_email: string | null, icanteen_password: string | null, token: string, connection: mysql.Connection) {
    this.email = email;
    this.password = password;
    this.icanteen_email = icanteen_email;
    this.icanteen_password = icanteen_password;
    this.token = token;
    this.page = null;
    this.connection = connection;
    this.last_active = null;
  }

  static async fromEmailPassword(email: string, password: string, connection: mysql.Connection): Promise<EP2User | boolean | null> {
    return new Promise((resolve) => {
      connection.query('SELECT * FROM users WHERE email = ?', [email], (err, res: Array<RowDataPacket>) => {
        if(err) {
          resolve(null)
        }
        if(res.length > 0) {
          const users = res.filter(u => u.password == password);
          if(users.length == 1) {
            const u = users[0]
            resolve(new EP2User(u.email, u.password, u.icanteen_email, u.icanteen_password, u.token, connection));
          } else {
            resolve(true);
          }
        } else {
          resolve(false);
        }
      })
    })
  }

  static async fromToken(token: string, connection: mysql.Connection): Promise<EP2User | boolean | null> {
    return new Promise((resolve) => {
      connection.query('SELECT * FROM users WHERE token = ?', [token], (err, res: Array<RowDataPacket>) => {
        if(err) {
          resolve(null)
        }
        if(res.length == 1) {
          const u = res[0];
          resolve(new EP2User(u.email, u.password, u.icanteen_email, u.icanteen_password, u.token, connection));
        } else if(res.length > 1) {
          resolve(true)
        } else {
          resolve(false)
        }
      })
    })
  }

  static async createFromEmailPassword(email: string, password: string, connection: mysql.Connection): Promise<EP2User> {
    return new Promise((resolve, reject) => {
      const token = randomAsciiString(32);
      connection.query('INSERT INTO users (email, password, token) VALUES (?, ?, ?)', [email, password, token], (err, _, __) =>{
        if(err) {
          reject(err)
        }
        resolve(new EP2User(email, password, null, null, token, connection));
      })
    })
  }

  async ICanteenCreatePage(browser: Browser, login = true, force = false) {
    if((this.icanteen_email != null && this.icanteen_password != null) || force) {
      this.page = await browser.newPage();
      await this.page.setViewport({width: 1080, height: 1024});
      await this.page.goto('https://stravovani.sspbrno.cz/login');
      if(login) {
        await this.ICanteenLogin();
      }
    }
  }

  async ICanteenLogin() {
    if(this.page == null || this.icanteen_email == null || this.icanteen_password == null) return;
    await this.page.type('#j_username', this.icanteen_email);
    await this.page.type('#j_password', this.icanteen_password);
    await this.page.click('.btn.btn-primary.btn-login');
  }

  async ICanteenLunchesForMonth() {
    if(this.page == null) return;
    if(this.page.url() != 'https://stravovani.sspbrno.cz/faces/secured/month.jsp') {
      await this.page.goto('https://stravovani.sspbrno.cz/faces/secured/month.jsp')
    }
    return await this.page.$$eval('#mainContext table tbody tr td form .orderContent', tds => tds.map((td) => {
      let returnArray = []
      for(let i of td.children[0].children) {
        returnArray.push({
          ordered: i.children[0].children[0].children[0].classList.contains('ordered'),
          can_order: !i.children[0].children[0].children[0].classList.contains('disabled'),
          item_name: i.children[0].children[1].textContent
        });
      }
      return {
        day: td.id.replace('orderContent', ''),
        lunches: returnArray
      };
    }))
  }
}